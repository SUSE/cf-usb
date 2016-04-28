package mgmt

import (
	"encoding/json"

	goerrors "errors"
	"fmt"
	"strings"

	"github.com/frodenas/brokerapi"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"

	uuid "github.com/satori/go.uuid"
)

func ConfigureInstanceAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger) {

	log := logger.Session("usb-mgmt-instance")

	api.CreateInstanceHandler = CreateInstanceHandlerFunc(func(params CreateInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("create-driver-instance")
		log.Info("request", lager.Data{"id": params.Instance.ID, "driver-instance-name": params.Instance.Name})

		if strings.ContainsAny(*params.Instance.Name, " ") {
			return &CreateInstanceInternalServerError{Payload: fmt.Sprintf("Driver instance name cannot contain spaces")}
		}

		var instance config.Instance

		instanceID := uuid.NewV4().String()

		instanceConfig, err := json.Marshal(params.Instance.Configuration)
		if err != nil {
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.TargetURL = params.Instance.TargetURL

		driverInstanceNameExist, err := configProvider.InstanceNameExists(*params.Instance.Name)
		if err != nil {
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		serviceNameExist := ccServiceBroker.CheckServiceNameExists(*params.Instance.Name)
		if driverInstanceNameExist || serviceNameExist {
			err := goerrors.New("A driver instance with the same name already exists")
			log.Error("check-driver-instance-name-exist", err)
			return &CreateInstanceConflict{}
		}
		instance.Name = *params.Instance.Name

		/*		driversPath, err := configProvider.GetDriversPath()
				if err != nil {
					return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
				}

						err = lib.Validate(instance, driversPath, existingDriver.DriverType, logger)
						if err != nil {
							log.Error("validation-failed", err)
							return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
						}
		*/
		err = configProvider.SetInstance(instanceID, instance)
		if err != nil {
			log.Error("set-driver-instance-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		var defaultDial config.Dial

		defaultDialID := uuid.NewV4().String()

		var plan brokerapi.ServicePlan

		plan.ID = uuid.NewV4().String()
		plan.Description = "default plan"
		plan.Name = "default"
		plan.Free = true

		var meta brokerapi.ServicePlanMetadata

		meta.DisplayName = "default plan"

		plan.Metadata = &meta

		defaultDial.Plan = plan
		defaultDialConfig := json.RawMessage([]byte("{}"))
		defaultDial.Configuration = &defaultDialConfig

		err = configProvider.SetDial(instanceID, defaultDialID, defaultDial)
		if err != nil {
			log.Error("set-dial-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		params.Instance.Dials = append(params.Instance.Dials, defaultDialID)
		params.Instance.ID = instanceID

		var service brokerapi.Service

		service.ID = uuid.NewV4().String()
		service.Name = *params.Instance.Name
		service.Description = "Default service"
		service.Tags = []string{*params.Instance.Name}
		service.Bindable = true

		err = configProvider.SetService(instanceID, service)
		if err != nil {
			log.Error("set-service-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		params.Instance.Service = service.ID

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			log.Error("load-configuration-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		if guid == "" {
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			log.Error("create-or-update-service-broker-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &CreateInstanceInternalServerError{Payload: err.Error()}
		}

		return &CreateInstanceCreated{Payload: params.Instance}
	})

	api.UpdateInstanceHandler = UpdateInstanceHandlerFunc(func(params UpdateInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("update-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.InstanceID})

		instanceInfo, _, err := configProvider.GetInstance(params.InstanceID)
		if err != nil {
			return &UpdateInstanceInternalServerError{Payload: err.Error()}
		}
		if instanceInfo == nil {
			return &UpdateInstanceNotFound{}
		}

		instance := *instanceInfo
		instanceConfig, err := json.Marshal(params.InstanceConfig.Configuration)
		if err != nil {
			return &UpdateInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		if instanceInfo.Name != *params.InstanceConfig.Name {
			driverInstanceNameExist, err := configProvider.InstanceNameExists(*params.InstanceConfig.Name)
			if err != nil {
				return &UpdateInstanceInternalServerError{Payload: err.Error()}
			}

			if driverInstanceNameExist {
				err := goerrors.New("A driver instance with the same name already exists")
				log.Error("check-driver-instance-name-exist", err)
				return &UpdateInstanceConflict{}
			}
		}
		instance.Name = *params.InstanceConfig.Name

		/*		existingDriver, err := configProvider.GetDriver(params.DriverConfig.DriverID)
				if err != nil {
					return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
				}

				driversPath, err := configProvider.GetDriversPath()
				if err != nil {
					return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
				}

						err = lib.Validate(instance, driversPath, existingDriver.DriverType, logger)
						if err != nil {
							log.Error("validation-failed", err)
							return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
						}
		*/
		err = configProvider.SetInstance(params.InstanceID, instance)
		if err != nil {
			return &UpdateInstanceInternalServerError{Payload: err.Error()}
		}

		return &UpdateInstanceOK{Payload: params.InstanceConfig}
	})

	api.DeleteInstanceHandler = DeleteInstanceHandlerFunc(func(params DeleteInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.InstanceID})

		instance, _, err := configProvider.GetInstance(params.InstanceID)
		if err != nil {
			return &DeleteInstanceInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &DeleteInstanceNotFound{}
		}
		if ccServiceBroker.CheckServiceInstancesExist(instance.Service.Name) == true {
			return &DeleteInstanceInternalServerError{Payload: fmt.Sprintf("Cannot delete instance '%s', it still has provisioned service instances", instance.Name)}
		}
		err = configProvider.DeleteInstance(params.InstanceID)
		if err != nil {
			return &DeleteInstanceInternalServerError{Payload: err.Error()}
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &DeleteInstanceInternalServerError{Payload: err.Error()}
		}

		instanceCount := 0
		for _, _ = range config.Instances {
			instanceCount++
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		if instanceCount == 0 {
			err := ccServiceBroker.Delete(brokerName)
			if err != nil {
				log.Error("delete-service-broker-failed", err)
				return &DeleteInstanceInternalServerError{Payload: err.Error()}
			}
		} else {
			guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
			if err != nil {
				log.Error("get-service-broker-failed", err)
				return &DeleteInstanceInternalServerError{Payload: err.Error()}
			}

			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
			if err != nil {
				log.Error("update-service-broker-failed", err)
				return &DeleteInstanceInternalServerError{Payload: err.Error()}
			}
		}

		return &DeleteInstanceNoContent{}
	})

	api.GetInstanceHandler = GetInstanceHandlerFunc(func(params GetInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("get-instance")
		log.Info("request", lager.Data{"instance-id": params.InstanceID})

		instance, _, err := configProvider.GetInstance(params.InstanceID)
		if err != nil {
			return &GetInstanceInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &GetInstanceNotFound{}
		}

		var conf map[string]interface{}

		err = json.Unmarshal(*instance.Configuration, &conf)
		if err != nil {
			return &GetInstanceInternalServerError{Payload: err.Error()}
		}

		var dials = make([]string, 0)

		for dialID, _ := range instance.Dials {
			dials = append(dials, dialID)
		}

		driverInstance := &genmodel.Instance{
			Configuration: conf,
			Dials:         dials,
			ID:            params.InstanceID,
			Name:          &instance.Name,
			Service:       instance.Service.ID,
		}

		return &GetInstanceOK{Payload: driverInstance}
	})

	api.GetInstancesHandler = GetInstancesHandlerFunc(func(principal interface{}) middleware.Responder {
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetInstancesInternalServerError{Payload: err.Error()}
		}
		var response []*genmodel.Instance
		for id, instance := range config.Instances {

			var dialIds []string
			for dialID, _ := range instance.Dials {
				dialIds = append(dialIds, dialID)
			}

			driverInstance := &genmodel.Instance{
				Configuration: instance.Configuration,
				TargetURL:     instance.TargetURL,
				Dials:         dialIds,
				ID:            id,
				Name:          &instance.Name,
				Service:       instance.Service.ID,
			}
			response = append(response, driverInstance)
		}
		return &GetInstancesOK{Payload: response}
	})
}
