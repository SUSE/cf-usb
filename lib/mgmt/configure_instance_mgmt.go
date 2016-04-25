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

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("create-driver-instance")
		log.Info("request", lager.Data{"driver-id": params.DriverInstance.DriverID, "driver-instance-name": params.DriverInstance.Name})

		if strings.ContainsAny(*params.DriverInstance.Name, " ") {
			return &CreateDriverInstanceInternalServerError{Payload: fmt.Sprintf("Driver instance name cannot contain spaces")}
		}

		existingDriver, err := configProvider.GetDriver(*params.DriverInstance.DriverID)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if existingDriver == nil {
			return &GetDriverNotFound{}
		}

		var instance config.DriverInstance

		instanceID := uuid.NewV4().String()

		instanceConfig, err := json.Marshal(params.DriverInstance.Configuration)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.TargetURL = params.DriverInstance.TargetURL

		driverInstanceNameExist, err := configProvider.DriverInstanceNameExists(*params.DriverInstance.Name)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		serviceNameExist := ccServiceBroker.CheckServiceNameExists(*params.DriverInstance.Name)
		if driverInstanceNameExist || serviceNameExist {
			err := goerrors.New("A driver instance with the same name already exists")
			log.Error("check-driver-instance-name-exist", err)
			return &CreateDriverInstanceConflict{}
		}
		instance.Name = *params.DriverInstance.Name

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
		err = configProvider.SetDriverInstance(*params.DriverInstance.DriverID, instanceID, instance)
		if err != nil {
			log.Error("set-driver-instance-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
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
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		params.DriverInstance.Dials = append(params.DriverInstance.Dials, defaultDialID)
		params.DriverInstance.ID = instanceID

		var service brokerapi.Service

		service.ID = uuid.NewV4().String()
		service.Name = *params.DriverInstance.Name
		service.Description = "Default service"
		service.Tags = []string{*params.DriverInstance.Name}
		service.Bindable = true

		err = configProvider.SetService(instanceID, service)
		if err != nil {
			log.Error("set-service-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		params.DriverInstance.Service = service.ID

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			log.Error("load-configuration-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		if guid == "" {
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			log.Error("create-or-update-service-broker-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		return &CreateDriverInstanceCreated{Payload: params.DriverInstance}
	})

	api.UpdateDriverInstanceHandler = UpdateDriverInstanceHandlerFunc(func(params UpdateDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("update-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		instanceInfo, _, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if instanceInfo == nil {
			return &UpdateDriverInstanceNotFound{}
		}

		instance := *instanceInfo
		instanceConfig, err := json.Marshal(params.DriverConfig.Configuration)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		if instanceInfo.Name != *params.DriverConfig.Name {
			driverInstanceNameExist, err := configProvider.DriverInstanceNameExists(*params.DriverConfig.Name)
			if err != nil {
				return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
			}

			if driverInstanceNameExist {
				err := goerrors.New("A driver instance with the same name already exists")
				log.Error("check-driver-instance-name-exist", err)
				return &UpdateDriverInstanceConflict{}
			}
		}
		instance.Name = *params.DriverConfig.Name

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
		err = configProvider.SetDriverInstance(*params.DriverConfig.DriverID, params.DriverInstanceID, instance)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		return &UpdateDriverInstanceOK{Payload: params.DriverConfig}
	})

	api.DeleteDriverInstanceHandler = DeleteDriverInstanceHandlerFunc(func(params DeleteDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		instance, _, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &DeleteDriverInstanceNotFound{}
		}
		if ccServiceBroker.CheckServiceInstancesExist(instance.Service.Name) == true {
			return &DeleteDriverInstanceInternalServerError{Payload: fmt.Sprintf("Cannot delete instance '%s', it still has provisioned service instances", instance.Name)}
		}
		err = configProvider.DeleteDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
		}

		instanceCount := 0
		for _, driver := range config.Drivers {
			for _, _ = range driver.DriverInstances {
				instanceCount++
			}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		if instanceCount == 0 {
			err := ccServiceBroker.Delete(brokerName)
			if err != nil {
				log.Error("delete-service-broker-failed", err)
				return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
			}
		} else {
			guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
			if err != nil {
				log.Error("get-service-broker-failed", err)
				return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
			}

			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
			if err != nil {
				log.Error("update-service-broker-failed", err)
				return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
			}
		}

		return &DeleteDriverInstanceNoContent{}
	})

	api.GetDriverInstanceHandler = GetDriverInstanceHandlerFunc(func(params GetDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		instance, driverID, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &GetDriverInstanceNotFound{}
		}

		var conf map[string]interface{}

		err = json.Unmarshal(*instance.Configuration, &conf)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}

		var dials = make([]string, 0)

		for dialID, _ := range instance.Dials {
			dials = append(dials, dialID)
		}

		driverInstance := &genmodel.DriverInstance{
			Configuration: conf,
			Dials:         dials,
			ID:            params.DriverInstanceID,
			Name:          &instance.Name,
			Service:       instance.Service.ID,
			DriverID:      &driverID,
		}

		return &GetDriverInstanceOK{Payload: driverInstance}
	})

}
