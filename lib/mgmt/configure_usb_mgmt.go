package mgmt

import (
	"encoding/json"

	"github.com/fatih/structs"
	"github.com/hpcloud/cf-usb/lib"

	"bufio"
	"crypto/sha1"
	"encoding/base64"
	goerrors "errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/frodenas/brokerapi"
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
	"github.com/xeipuuv/gojsonschema"

	uuid "github.com/satori/go.uuid"
)

// This file is safe to edit. Once it exists it will not be overwritten

const defaultBrokerName string = "usb"

func ConfigureAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface, configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface, logger lager.Logger) {
	log := logger.Session("usb-mgmt")

	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.AuthorizationAuth = func(token string) (interface{}, error) {
		err := auth.IsAuthenticated(token)

		if err != nil {
			return nil, err
		}

		log.Debug("authentication-succeeded")

		return true, nil
	}

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("update-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		exists, err := configProvider.DriverExists(params.DriverID)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}
		if exists == false {
			log.Debug("update-driver-does-not-exist", lager.Data{"driver-id-does-not-exit": params.DriverID})
			return &UpdateDriverNotFound{}
		}

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}

		driver.DriverName = params.Driver.Name

		err = configProvider.SetDriver(*params.Driver.ID, *driver)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}

		return &UpdateDriverOK{Payload: params.Driver}
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("upload-driver-bits")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &UploadDriverNotFound{}
		}

		driverType := driver.DriverType

		driverPath := os.Getenv("USB_DRIVER_PATH")
		if driverPath == "" {
			driverPath = "drivers"
		}
		driverPath = filepath.Join(driverPath, driverType)
		if runtime.GOOS == "windows" {
			driverPath = driverPath + ".exe"
		}

		defer params.File.Data.Close()

		f, err := os.OpenFile(driverPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}

		defer f.Close()

		reader := bufio.NewReader(params.File.Data)

		sha1 := sha1.New()
		_, err = io.Copy(f, io.TeeReader(reader, sha1))
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}

		sha := base64.StdEncoding.EncodeToString(sha1.Sum(nil))
		if sha != params.Sha {
			f.Close()
			os.Remove(driverPath)

			return &UploadDriverInternalServerError{Payload: fmt.Sprintf("Checksum mismatch. Actual file SHA1 checksum: %s. Expected file SHA1 checksum: %s", sha, params.Sha)}
		}

		return &UploadDriverOK{}
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func(params DeleteServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-service-plan")
		log.Info("request", lager.Data{"plan-id": params.PlanID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &DeleteServicePlanInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &DeleteServicePlanInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			log.Error("update-service-broker-failed", err)
			return &DeleteServicePlanInternalServerError{Payload: err.Error()}
		}

		//TODO improve this
		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for dialID, dial := range instance.Dials {
					if dial.Plan.ID == params.PlanID {
						err = configProvider.DeleteDial(dialID)
						if err != nil {
							return &DeleteServicePlanInternalServerError{Payload: err.Error()}
						}
						return &DeleteServicePlanNoContent{}
					}
				}

				err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
				if err != nil {
					log.Error("enable-service-access-failed", err)
					return &DeleteServicePlanInternalServerError{Payload: err.Error()}
				}
			}
		}

		return &DeleteServicePlanNotFound{}
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-schema")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		path, err := configProvider.GetDriversPath()
		if err != nil {
			return &GetDriverSchemaInternalServerError{Payload: err.Error()}
		}
		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverSchemaNotFound{}
		}
		schema, err := lib.GetConfigSchema(path, driver.DriverType, logger)
		if err != nil {
			return &GetDriverSchemaInternalServerError{Payload: err.Error()}
		}

		return &GetDriverSchemaOK{Payload: genmodel.DriverSchema(schema)}
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		d, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverNotFound{}
		}

		var instances = make([]string, 0)
		for instanceID, _ := range d.DriverInstances {
			instances = append(instances, instanceID)
		}

		driver := &genmodel.Driver{
			DriverType:      d.DriverType,
			DriverInstances: instances,
			ID:              &params.DriverID,
		}

		return &GetDriverOK{Payload: driver}
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("get-drivers")
		log.Info("request")

		var drivers = make([]*genmodel.Driver, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDriversInternalServerError{Payload: err.Error()}
		}

		for dId, d := range config.Drivers {
			var instances = make([]string, 0)
			for instanceID, _ := range d.DriverInstances {
				instances = append(instances, instanceID)
			}

			var driverID string
			driverID = dId

			driver := &genmodel.Driver{
				ID:              &driverID,
				DriverType:      d.DriverType,
				Name:            d.DriverName,
				DriverInstances: instances,
			}

			drivers = append(drivers, driver)
		}

		log.Debug("", lager.Data{"drivers-found": len(drivers)})

		return &GetDriversOK{Payload: drivers}
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("create-driver")
		log.Info("request", lager.Data{"driver-name": params.Driver.Name, "driver-type": params.Driver.DriverType})

		exist, err := configProvider.DriverTypeExists(params.Driver.DriverType)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}
		if exist {
			return &CreateDriverConflict{}
		}

		var driver config.Driver

		driver.DriverType = params.Driver.DriverType
		driver.DriverName = params.Driver.Name

		driverID := uuid.NewV4().String()

		err = configProvider.SetDriver(driverID, driver)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}

		params.Driver.ID = &driverID

		return &CreateDriverCreated{Payload: params.Driver}
	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-instances")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		var driverInstances = make([]*genmodel.DriverInstance, 0)

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}

		for diID, di := range driver.DriverInstances {
			var dials = make([]string, 0)
			for dialID, _ := range di.Dials {
				dials = append(dials, dialID)
			}

			var driverInstanceID string
			driverInstanceID = diID

			driverInstance := &genmodel.DriverInstance{
				Configuration: di.Configuration,
				Dials:         dials,
				DriverID:      params.DriverID,
				ID:            &driverInstanceID,
				Name:          di.Name,
				Service:       &di.Service.ID,
			}

			driverInstances = append(driverInstances, driverInstance)
		}

		log.Debug("", lager.Data{"driver-instances-found": len(driverInstances)})

		return &GetDriverInstancesOK{Payload: driverInstances}
	})

	api.PingDriverInstanceHandler = PingDriverInstanceHandlerFunc(func(params PingDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("ping-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		configuration, err := configProvider.LoadConfiguration()
		if err != nil {
			return &PingDriverInstanceInternalServerError{Payload: err.Error()}
		}

		for _, driver := range configuration.Drivers {
			for instanceID, instance := range driver.DriverInstances {
				if instanceID == params.DriverInstanceID {
					result, err := lib.Ping(instance.Configuration, configuration.DriversPath, driver.DriverType)
					if err != nil {
						return &PingDriverInstanceInternalServerError{Payload: err.Error()}
					}
					if result == true {
						return &PingDriverInstanceOK{}
					} else {
						return &PingDriverInstanceNotFound{}
					}
				}
			}
		}
		return &PingDriverInstanceNotFound{}
	})

	api.GetServicePlanHandler = GetServicePlanHandlerFunc(func(params GetServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service-plan")
		log.Info("request", lager.Data{"plan-id": params.PlanID})

		planInfo, dialID, _, err := configProvider.GetPlan(params.PlanID)
		if err != nil {
			return &GetServicePlanNotFound{}
		}

		plan := &genmodel.Plan{
			Name:        planInfo.Name,
			ID:          &planInfo.ID,
			DialID:      dialID,
			Description: &planInfo.Description,
			Free:        &planInfo.Free,
		}

		return &GetServicePlanOK{Payload: plan}

	})

	api.CreateDialHandler = CreateDialHandlerFunc(func(params CreateDialParams, principal interface{}) middleware.Responder {
		log := log.Session("create-dial")
		log.Info("request", lager.Data{"driver-instance-id": params.Dial.DriverInstanceID})

		var dial config.Dial

		dialID := uuid.NewV4().String()
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		err = configProvider.SetDial(params.Dial.DriverInstanceID, dialID, dial)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		params.Dial.ID = &dialID

		return &CreateDialCreated{Payload: params.Dial}
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams, principal interface{}) middleware.Responder {
		log := log.Session("update-dial")
		log.Info("request", lager.Data{"driver-instance-id": params.Dial.DriverInstanceID})

		if params.DialID == "" {
			return &UpdateDialInternalServerError{Payload: "Empty dial id parameter provided"}
		}

		if params.Dial == nil {
			return &UpdateDialInternalServerError{Payload: "Empty dial provided"}
		}

		dialID := params.DialID

		dial, err := configProvider.GetDial(dialID)
		if err != nil {
			log.Error("update-dial", err, lager.Data{"dial-id": dialID})
			return &UpdateDialNotFound{}
		}

		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		_, parentId, err := configProvider.GetDriverInstance(params.Dial.DriverInstanceID)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		driver, err := configProvider.GetDriver(parentId)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}
		path, err := configProvider.GetDriversPath()
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		schema, err := lib.GetDailsSchema(path, driver.DriverType, logger)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}
		dialsSchemaLoader := gojsonschema.NewStringLoader(schema)
		dialLoader := gojsonschema.NewGoLoader(dial.Configuration)
		result, err := gojsonschema.Validate(dialsSchemaLoader, dialLoader)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}
		if !result.Valid() {
			err = goerrors.New("Invalid dial configuration")
			logger.Error("update-dial-validate-schema", err, lager.Data{"Errors": result.Errors()})
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		err = configProvider.SetDial(params.Dial.DriverInstanceID, dialID, *dial)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		return &UpdateDialOK{Payload: params.Dial}
	})

	api.GetDialSchemaHandler = GetDialSchemaHandlerFunc(func(params GetDialSchemaParams, principal interface{}) middleware.Responder {
		log := log.Session("get-dial-schema")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		path, err := configProvider.GetDriversPath()
		if err != nil {
			return &GetDialSchemaInternalServerError{Payload: err.Error()}
		}
		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDialSchemaNotFound{}
		}
		schema, err := lib.GetDailsSchema(path, driver.DriverType, logger)
		if err != nil {
			return &GetDialSchemaInternalServerError{Payload: err.Error()}
		}
		return &GetDialSchemaOK{Payload: genmodel.DialSchema(schema)}
	})

	api.GetServiceByInstanceIDHandler = GetServiceByInstanceIDHandlerFunc(func(params GetServiceByInstanceIDParams, principal interface{}) middleware.Responder {
		log := log.Session("get-services")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		di, _, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &GetServiceByInstanceIDInternalServerError{Payload: err.Error()}
		}

		service := &genmodel.Service{
			ID:               &di.Service.ID,
			DriverInstanceID: params.DriverInstanceID,
			Bindable:         &di.Service.Bindable,
			Name:             di.Service.Name,
			Tags:             di.Service.Tags,
		}

		if di.Service.Metadata != nil {
			service.Metadata = structs.Map(*di.Service.Metadata)
		}

		return &GetServiceByInstanceIDOK{Payload: service}
	})

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("create-driver-instance")
		log.Info("request", lager.Data{"driver-id": params.DriverInstance.DriverID, "driver-instance-name": params.DriverInstance.Name})

		existingDriver, err := configProvider.GetDriver(params.DriverInstance.DriverID)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		var instance config.DriverInstance

		instanceID := uuid.NewV4().String()

		instanceConfig, err := json.Marshal(params.DriverInstance.Configuration)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration

		driverInstanceNameExist, err := configProvider.DriverInstanceNameExists(params.DriverInstance.Name)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		serviceNameExist := ccServiceBroker.CheckServiceNameExists(params.DriverInstance.Name)
		if driverInstanceNameExist || serviceNameExist {
			err := goerrors.New("A driver instance with the same name already exists")
			log.Error("check-driver-instance-name-exist", err)
			return &CreateDriverInstanceConflict{}
		}
		instance.Name = params.DriverInstance.Name

		driversPath, err := configProvider.GetDriversPath()
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		err = lib.Validate(instance, driversPath, existingDriver.DriverType, logger)
		if err != nil {
			log.Error("validation-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		err = configProvider.SetDriverInstance(params.DriverInstance.DriverID, instanceID, instance)
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
		params.DriverInstance.ID = &instanceID

		var service brokerapi.Service

		service.ID = uuid.NewV4().String()
		service.Name = params.DriverInstance.Name
		service.Description = "Default service"
		service.Tags = []string{params.DriverInstance.Name}
		service.Bindable = true

		err = configProvider.SetService(instanceID, service)
		if err != nil {
			log.Error("set-service-failed", err)
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		params.DriverInstance.Service = &service.ID

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
			return &UpdateDriverInstanceNotFound{}
		}

		instance := *instanceInfo
		instanceConfig, err := json.Marshal(params.DriverConfig.Configuration)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		if instanceInfo.Name != params.DriverConfig.Name {
			driverInstanceNameExist, err := configProvider.DriverInstanceNameExists(params.DriverConfig.Name)
			if err != nil {
				return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
			}

			if driverInstanceNameExist {
				err := goerrors.New("A driver instance with the same name already exists")
				log.Error("check-driver-instance-name-exist", err)
				return &UpdateDriverInstanceConflict{}
			}
		}
		instance.Name = params.DriverConfig.Name

		existingDriver, err := configProvider.GetDriver(params.DriverConfig.DriverID)
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

		err = configProvider.SetDriverInstance(params.DriverConfig.DriverID, params.DriverInstanceID, instance)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		return &UpdateDriverInstanceOK{Payload: params.DriverConfig}
	})

	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		_, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &DeleteDriverNotFound{}
		}

		err = configProvider.DeleteDriver(params.DriverID)
		if err != nil {
			return &DeleteDriverInternalServerError{Payload: err.Error()}
		}

		return &DeleteDriverNoContent{}
	})

	api.DeleteDriverInstanceHandler = DeleteDriverInstanceHandlerFunc(func(params DeleteDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		_, _, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &DeleteDriverInstanceNotFound{}
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

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams, principal interface{}) middleware.Responder {
		log := log.Session("get-dials")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		var dials = make([]*genmodel.Dial, 0)
		if *params.DriverInstanceID == "" {
			return &GetAllDialsInternalServerError{Payload: "Empty driver instance id in get all dials"}
		}
		instanceInfo, err := configProvider.LoadDriverInstance(*params.DriverInstanceID)
		if err != nil {
			return &GetAllDialsInternalServerError{Payload: err.Error()}
		}

		for diaID, dia := range instanceInfo.Dials {

			dial := &genmodel.Dial{
				Configuration:    dia.Configuration,
				DriverInstanceID: *params.DriverInstanceID,
				ID:               &diaID,
				Plan:             &dia.Plan.ID,
			}

			dials = append(dials, dial)
		}

		log.Debug("", lager.Data{"dials-found": len(dials)})

		return &GetAllDialsOK{Payload: dials}
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("update-service-plan")
		log.Info("request", lager.Data{"plan-id": params.PlanID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}

		_, _, _, err = configProvider.GetPlan(params.PlanID)
		if err != nil {
			return &UpdateServicePlanNotFound{}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			log.Error("update-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for instanceID, instance := range driver.DriverInstances {
				for dialID, dial := range instance.Dials {
					if dialID == params.Plan.DialID {
						if dial.Plan.ID == params.PlanID {
							var plan brokerapi.ServicePlan
							var meta brokerapi.ServicePlanMetadata

							plan.Description = *params.Plan.Description
							plan.ID = *params.Plan.ID
							plan.Name = params.Plan.Name
							plan.Free = *params.Plan.Free

							meta.DisplayName = params.Plan.Name
							plan.Metadata = &meta
							dial.Plan = plan
							err = configProvider.SetDial(instanceID, dialID, dial)

							if err != nil {
								return &UpdateServicePlanInternalServerError{Payload: err.Error()}
							}
							return &UpdateServicePlanOK{Payload: params.Plan}
						}
					}
				}

				err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
				if err != nil {
					log.Error("enable-service-access-failed", err)
					return &UpdateServiceInternalServerError{Payload: err.Error()}
				}
			}
		}

		return &UpdateServicePlanNotFound{}
	})

	api.GetServiceHandler = GetServiceHandlerFunc(func(params GetServiceParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service")
		log.Info("request", lager.Data{"service-id": params.ServiceID})

		serviceInfo, instanceID, err := configProvider.GetService(params.ServiceID)
		if err != nil {
			return &GetServiceNotFound{}
		}

		svc := &genmodel.Service{
			Bindable:         &serviceInfo.Bindable,
			DriverInstanceID: instanceID,
			ID:               &serviceInfo.ID,
			Name:             serviceInfo.Name,
			Description:      &serviceInfo.Description,
			Tags:             serviceInfo.Tags,
		}

		if serviceInfo.Metadata != nil {
			svc.Metadata = structs.Map(*serviceInfo.Metadata)
		}

		return &GetServiceOK{Payload: svc}
	})

	api.DeleteDialHandler = DeleteDialHandlerFunc(func(params DeleteDialParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-dial")
		log.Info("request", lager.Data{"dial-id": params.DialID})

		_, err := configProvider.GetDial(params.DialID)
		if err != nil {
			return &DeleteDialNotFound{}
		}

		err = configProvider.DeleteDial(params.DialID)
		if err != nil {
			return &DeleteDialInternalServerError{Payload: err.Error()}
		}
		return &DeleteDialNoContent{}
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("get-info")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetInfoInternalServerError{
				Payload: err.Error(),
			}
		}

		info := &genmodel.Info{
			Version: config.APIVersion,
		}

		return &GetInfoOK{
			Payload: info,
		}
	})

	api.GetServicePlansHandler = GetServicePlansHandlerFunc(func(params GetServicePlansParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service-plans")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		if *params.DriverInstanceID == "" {
			return &GetServicePlansInternalServerError{Payload: "Empty driver instance id in get service plans handler"}
		}
		var servicePlans = make([]*genmodel.Plan, 0)

		instanceInfo, err := configProvider.LoadDriverInstance(*params.DriverInstanceID)
		if err != nil {
			return &GetServicePlansInternalServerError{Payload: err.Error()}
		}

		for diaID, dia := range instanceInfo.Dials {
			plan := &genmodel.Plan{
				Name:        dia.Plan.Name,
				ID:          &dia.Plan.ID,
				DialID:      diaID,
				Description: &dia.Plan.Description,
				Free:        &dia.Plan.Free,
			}

			servicePlans = append(servicePlans, plan)
		}

		log.Debug("", lager.Data{"service-plans-found": len(servicePlans)})

		return &GetServicePlansOK{Payload: servicePlans}
	})

	api.GetDialHandler = GetDialHandlerFunc(func(params GetDialParams, principal interface{}) middleware.Responder {
		log := log.Session("get-dial")
		log.Info("request", lager.Data{"dial-id": params.DialID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDialInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			for diID, di := range d.DriverInstances {
				for diaID, dia := range di.Dials {
					if diaID == params.DialID {

						var conf map[string]interface{}
						err := json.Unmarshal(*dia.Configuration, &conf)
						if err != nil {
							return &GetDialInternalServerError{Payload: err.Error()}
						}

						dial := &genmodel.Dial{
							Configuration:    conf,
							DriverInstanceID: diID,
							ID:               &diaID,
							Plan:             &dia.Plan.ID,
						}

						return &GetDialOK{Payload: dial}
					}
				}
			}
		}

		return &GetDialNotFound{}
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams, principal interface{}) middleware.Responder {
		log := log.Session("update-service")
		log.Info("request", lager.Data{"service-id": params.ServiceID})

		if params.Service == nil {
			return &UpdateServiceInternalServerError{Payload: "Service information cannot be nil"}
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		service, instanceid, err := configProvider.GetService(params.ServiceID)
		if err != nil {
			return &UpdateServiceNotFound{}
		}
		if service == nil {
			return &UpdateServiceInternalServerError{Payload: fmt.Sprintf("Service-id %s not found", params.ServiceID)}
		}

		if params.Service.Bindable != nil {
			service.Bindable = *params.Service.Bindable
		}
		if params.Service.Description != nil {
			service.Description = *params.Service.Description
		}
		if params.Service.Name != "" {
			service.Name = params.Service.Name
		}
		if len(params.Service.Tags) > 0 {
			service.Tags = params.Service.Tags
		}

		exists := ccServiceBroker.CheckServiceNameExists(params.Service.Name)

		if exists == true {
			err := goerrors.New("Service update name parameter validation failed - duplicate naming eror")
			log.Error("update-service-name-validation", err, lager.Data{"Name validation failed for name": params.Service.Name})
			return &UpdateServiceConflict{}
		}

		err = configProvider.SetService(instanceid, *service)
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			log.Error("update-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		return &UpdateServiceOK{Payload: params.Service}
	})

	api.GetDriverInstanceHandler = GetDriverInstanceHandlerFunc(func(params GetDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

		instance, _, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
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
			ID:            &params.DriverInstanceID,
			Name:          instance.Name,
			Service:       &instance.Service.ID,
		}

		return &GetDriverInstanceOK{Payload: driverInstance}
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams, principal interface{}) middleware.Responder {
		log := log.Session("create-service-plan")
		log.Info("request", lager.Data{"dial-id": params.Plan.DialID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &CreateServicePlanInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			log.Error("update-service-broker-failed", err)
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for instanceID, instance := range driver.DriverInstances {
				for dialID, dial := range instance.Dials {
					if dialID == params.Plan.DialID {
						err = configProvider.DeleteDial(dialID)
						if err != nil {
							return &CreateServicePlanInternalServerError{Payload: err.Error()}
						}

						var plan brokerapi.ServicePlan
						var meta brokerapi.ServicePlanMetadata

						plan.Description = *params.Plan.Description
						plan.ID = uuid.NewV4().String()
						plan.Name = params.Plan.Name
						plan.Free = *params.Plan.Free

						meta.DisplayName = params.Plan.Name

						plan.Metadata = &meta

						dial.Plan = plan

						err = configProvider.SetDial(instanceID, dialID, dial)
						if err != nil {
							return &CreateServicePlanInternalServerError{Payload: err.Error()}
						}

						params.Plan.ID = &plan.ID

						return &CreateServicePlanCreated{Payload: params.Plan}
					}
				}

				err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
				if err != nil {
					log.Error("enable-service-access-failed", err)
					return &UpdateServiceInternalServerError{Payload: err.Error()}
				}
			}
		}

		return &CreateServicePlanInternalServerError{Payload: fmt.Sprintf("Dial %s not found", params.Plan.DialID)}
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("update-catalog")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		if guid == "" {
			return &UpdateCatalogInternalServerError{Payload: fmt.Sprintf("Broker %s guid not found", brokerName)}
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
			if err != nil {
				return &UpdateCatalogInternalServerError{Payload: err.Error()}
			}
		}
		return &UpdateCatalogOK{}
	})
}
