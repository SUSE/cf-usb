package mgmt

import (
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/hpcloud/cf-usb/lib"

	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// This file is safe to edit. Once it exists it will not be overwritten

const brokerName string = "usb"

func ConfigureAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface, configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface, logger lager.Logger) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.AuthorizationAuth = func(token string) (interface{}, error) {
		err := auth.IsAuthenticated(token)

		if err != nil {
			return nil, err
		}

		return true, nil
	}

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams, principal interface{}) middleware.Responder {
		_, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &UpdateDriverNotFound{}
		}

		var driver config.Driver

		driver.DriverType = params.Driver.DriverType
		driver.ID = params.Driver.ID

		err = configProvider.SetDriver(driver)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}

		return &UpdateDriverOK{Payload: params.Driver}
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams, principal interface{}) middleware.Responder {
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

			return &UploadDriverInternalServerError{Payload: fmt.Sprintf("Checksum mismatch. Expected: %s, got %s", params.Sha, sha)}
		}

		return &UploadDriverOK{}
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func(params DeleteServicePlanParams, principal interface{}) middleware.Responder {
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &DeleteServicePlanInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.Plan.ID == params.PlanID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return &DeleteServicePlanInternalServerError{Payload: err.Error()}
						}
						return &DeleteServicePlanNoContent{}
					}
				}
			}
		}

		return &DeleteServicePlanNotFound{}
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation getDriverSchema has not yet been implemented")
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams, principal interface{}) middleware.Responder {
		d, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverNotFound{}
		}

		var instances = make([]string, 0)
		for _, instance := range d.DriverInstances {
			instances = append(instances, instance.ID)
		}

		driver := &genmodel.Driver{
			ID:              d.ID,
			DriverType:      d.DriverType,
			DriverInstances: instances,
		}
		return &GetDriverOK{Payload: driver}
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(principal interface{}) middleware.Responder {

		var drivers = make([]*genmodel.Driver, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDriversInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			var instances = make([]string, 0)
			for _, instance := range d.DriverInstances {
				instances = append(instances, instance.ID)
			}

			driver := &genmodel.Driver{
				ID:              d.ID,
				DriverType:      d.DriverType,
				DriverInstances: instances,
			}

			drivers = append(drivers, driver)
		}

		return &GetDriversOK{Payload: drivers}
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams, principal interface{}) middleware.Responder {

		exist, err := configProvider.DriverTypeExists(params.Driver.DriverType)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}
		if exist {
			return &CreateDriverConflict{}
		}

		var driver config.Driver

		driver.DriverType = params.Driver.DriverType

		driver.ID = uuid.NewV4().String()

		err = configProvider.SetDriver(driver)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}
		params.Driver.ID = driver.ID
		return &CreateDriverCreated{Payload: params.Driver}
	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams, principal interface{}) middleware.Responder {

		var driverInstances = make([]*genmodel.DriverInstance, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			if params.DriverID != "" && d.ID != params.DriverID {
				continue
			}

			for _, di := range d.DriverInstances {
				var dials = make([]string, 0)
				for _, dial := range di.Dials {
					dials = append(dials, dial.ID)
				}

				var conf map[string]interface{}
				err := json.Unmarshal(*di.Configuration, &conf)
				if err != nil {
					return &GetDriverInstanceInternalServerError{Payload: err.Error()}
				}

				driverInstance := &genmodel.DriverInstance{
					Configuration: conf,
					Dials:         dials,
					ID:            di.ID,
					DriverID:      d.ID,
					Name:          di.Name,
					Service:       di.Service.ID,
				}

				driverInstances = append(driverInstances, driverInstance)
			}
		}

		return &GetDriverInstancesOK{Payload: driverInstances}
	})

	api.PingDriverInstanceHandler = PingDriverInstanceHandlerFunc(func(params PingDriverInstanceParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation pingDriverInstance has not yet been implemented")
	})

	api.GetServicePlanHandler = GetServicePlanHandlerFunc(func(params GetServicePlanParams, principal interface{}) middleware.Responder {
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetServicePlanInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				for _, dial := range di.Dials {
					if dial.Plan.ID == params.PlanID {
						plan := &genmodel.Plan{
							Name:        dial.Plan.Name,
							ID:          dial.Plan.ID,
							DialID:      dial.ID,
							Description: dial.Plan.Description,
							// TODO add free
						}

						return &GetServicePlanOK{Payload: plan}
					}
				}
			}
		}

		return &GetServicePlanNotFound{}
	})

	api.CreateDialHandler = CreateDialHandlerFunc(func(params CreateDialParams, principal interface{}) middleware.Responder {

		var dial config.Dial

		dial.ID = uuid.NewV4().String()
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}
		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		err = configProvider.SetDial(params.Dial.DriverInstanceID, dial)
		if err != nil {
			return &CreateDialInternalServerError{Payload: err.Error()}
		}

		params.Dial.ID = dial.ID

		return &CreateDialCreated{Payload: params.Dial}
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams, principal interface{}) middleware.Responder {
		var dial config.Dial

		dial.ID = params.Dial.ID
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration
		dial.Plan = brokerapi.ServicePlan{ID: params.Dial.Plan}
		err = configProvider.SetDial(params.Dial.DriverInstanceID, dial)
		if err != nil {
			return &UpdateDialInternalServerError{Payload: err.Error()}
		}
		return &UpdateDialOK{Payload: params.Dial}
	})

	api.GetDialSchemaHandler = GetDialSchemaHandlerFunc(func(params GetDialSchemaParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation getDialSchema has not yet been implemented")
	})

	api.GetServicesHandler = GetServicesHandlerFunc(func(params GetServicesParams, principal interface{}) middleware.Responder {
		var services = make([]*genmodel.Service, 0)

		di, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &GetServicesInternalServerError{Payload: err.Error()}
		}

		service := &genmodel.Service{
			ID:               di.Service.ID,
			DriverInstanceID: di.ID,
			Bindable:         di.Service.Bindable,
			Name:             di.Service.Name,
			Tags:             di.Service.Tags,
		}

		if di.Service.Metadata != nil {
			service.Metadata = structs.Map(*di.Service.Metadata)
		}

		services = append(services, service)

		return &GetServicesOK{Payload: services}
	})

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams, principal interface{}) middleware.Responder {

		existingDriver, err := configProvider.GetDriver(params.DriverInstance.DriverID)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		var instance config.DriverInstance
		instance.ID = uuid.NewV4().String()
		instanceConfig, err := json.Marshal(params.DriverInstance.Configuration)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.Name = params.DriverInstance.Name

		err = lib.Validate(instance, existingDriver.DriverType, logger)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		err = configProvider.SetDriverInstance(params.DriverInstance.DriverID, instance)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		var defaultDial config.Dial

		defaultDial.ID = uuid.NewV4().String()

		var plan brokerapi.ServicePlan
		plan.ID = uuid.NewV4().String()
		plan.Description = "default plan"
		plan.Name = "default"

		var meta brokerapi.ServicePlanMetadata
		meta.DisplayName = "default plan"

		plan.Metadata = &meta

		defaultDial.Plan = plan
		defaultDialConfig := json.RawMessage([]byte("{}"))
		defaultDial.Configuration = &defaultDialConfig

		err = configProvider.SetDial(instance.ID, defaultDial)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		params.DriverInstance.Dials = append(params.DriverInstance.Dials, defaultDial.ID)
		params.DriverInstance.ID = instance.ID

		var service brokerapi.Service
		service.ID = uuid.NewV4().String()
		service.Name = fmt.Sprintf("%s-%s", params.DriverInstance.Name, GetRandomString(5))
		service.Description = "Default service"
		service.Tags = []string{params.DriverInstance.Name}
		service.Bindable = true

		err = configProvider.SetService(instance.ID, service)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		params.DriverInstance.Service = service.ID
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		logger.Info("create-instance", lager.Data{"get-broker": brokerName})
		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)

		logger.Info("create-instance", lager.Data{"broker-guid": guid})
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		if guid == "" {

			logger.Info("create-instance", lager.Data{"service-broker-create": service.Name})
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {

			logger.Info("create-instance", lager.Data{"service-broker-update": service.Name})
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		logger.Info("create-instance", lager.Data{"enable-access": service.Name})
		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			return &CreateDriverInstanceInternalServerError{Payload: err.Error()}
		}

		return &CreateDriverInstanceCreated{Payload: params.DriverInstance}
	})

	api.UpdateDriverInstanceHandler = UpdateDriverInstanceHandlerFunc(func(params UpdateDriverInstanceParams, principal interface{}) middleware.Responder {
		instance, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &UpdateDriverInstanceNotFound{}
		}

		instanceConfig, err := json.Marshal(params.DriverConfig.Configuration)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.Name = params.DriverConfig.Name

		err = configProvider.SetDriverInstance(params.DriverConfig.DriverID, instance)
		if err != nil {
			return &UpdateDriverInstanceInternalServerError{Payload: err.Error()}
		}
		return &UpdateDriverInstanceOK{Payload: params.DriverConfig}
	})

	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams, principal interface{}) middleware.Responder {
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
		_, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &DeleteDriverInstanceNotFound{}
		}

		err = configProvider.DeleteDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &DeleteDriverInstanceInternalServerError{Payload: err.Error()}
		}

		return &DeleteDriverInstanceNoContent{}
	})

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams, principal interface{}) middleware.Responder {
		var dials = make([]*genmodel.Dial, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDialSchemaInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID {
					continue
				}

				for _, dia := range di.Dials {

					var conf map[string]interface{}
					err := json.Unmarshal(*dia.Configuration, &conf)
					if err != nil {
						return &GetDialSchemaInternalServerError{Payload: err.Error()}
					}

					dial := &genmodel.Dial{
						Configuration:    conf,
						DriverInstanceID: di.ID,
						ID:               dia.ID,
						Plan:             dia.Plan.ID,
					}

					dials = append(dials, dial)
				}
			}
		}

		return &GetAllDialsOK{Payload: dials}
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams, principal interface{}) middleware.Responder {

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServicePlanInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.ID == params.Plan.DialID {
						if dial.Plan.ID == params.PlanID {
							var plan brokerapi.ServicePlan
							var meta brokerapi.ServicePlanMetadata

							plan.Description = params.Plan.Description
							plan.ID = params.Plan.ID
							plan.Name = params.Plan.Name

							meta.DisplayName = params.Plan.Name
							plan.Metadata = &meta
							dial.Plan = plan
							err = configProvider.SetDial(instance.ID, dial)

							if err != nil {
								return &UpdateServicePlanInternalServerError{Payload: err.Error()}
							}
							return &UpdateServicePlanOK{Payload: params.Plan}
						}
					} else {
						return &UpdateServicePlanNotFound{}
					}
				}
			}
		}

		return &UpdateServicePlanNotFound{}
	})

	api.GetServiceHandler = GetServiceHandlerFunc(func(params GetServiceParams, principal interface{}) middleware.Responder {

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetServiceInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				if instance.Service.ID == params.ServiceID {
					svc := &genmodel.Service{
						Bindable:         instance.Service.Bindable,
						DriverInstanceID: instance.ID,
						ID:               instance.Service.ID,
						Name:             instance.Service.Name,
						Description:      instance.Service.Description,
						Tags:             instance.Service.Tags,
					}

					if instance.Service.Metadata != nil {
						svc.Metadata = structs.Map(*instance.Service.Metadata)
					}

					return &GetServiceOK{Payload: svc}
				}
			}
		}

		return &GetServiceNotFound{}
	})

	api.DeleteDialHandler = DeleteDialHandlerFunc(func(params DeleteDialParams, principal interface{}) middleware.Responder {

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &DeleteDialInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.ID == params.DialID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return &DeleteDialInternalServerError{Payload: err.Error()}
						}
						return &DeleteDialNoContent{}
					}
				}
			}
		}
		return &DeleteDialNotFound{}
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func(principal interface{}) middleware.Responder {
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
		var servicePlans = make([]*genmodel.Plan, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetServicePlansInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID {
					continue
				}

				for _, dia := range di.Dials {
					plan := &genmodel.Plan{
						Name:        dia.Plan.Name,
						ID:          dia.Plan.ID,
						DialID:      dia.ID,
						Description: dia.Plan.Description,
						// TODO add free
					}

					servicePlans = append(servicePlans, plan)
				}
			}
		}

		return &GetServicePlansOK{Payload: servicePlans}
	})

	api.GetDialHandler = GetDialHandlerFunc(func(params GetDialParams, principal interface{}) middleware.Responder {
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDialInternalServerError{Payload: err.Error()}
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				for _, dia := range di.Dials {
					if dia.ID == params.DialID {

						var conf map[string]interface{}
						err := json.Unmarshal(*dia.Configuration, &conf)
						if err != nil {
							return &GetDialInternalServerError{Payload: err.Error()}
						}

						dial := &genmodel.Dial{
							Configuration:    conf,
							DriverInstanceID: di.ID,
							ID:               dia.ID,
							Plan:             dia.Plan.ID,
						}

						return &GetDialOK{Payload: dial}
					}
				}
			}
		}

		return &GetDialNotFound{}
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams, principal interface{}) middleware.Responder {
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		var service brokerapi.Service

		service.Bindable = params.Service.Bindable
		service.Description = params.Service.Description
		service.ID = params.Service.ID
		service.Name = params.Service.Name
		service.Tags = params.Service.Tags

		err = configProvider.SetService(params.Service.DriverInstanceID, service)
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		logger.Info("update-service", lager.Data{"get-broker": brokerName})
		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)

		logger.Info("update-service", lager.Data{"broker-guid": guid})
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		logger.Info("update-service", lager.Data{"service-broker-update": service.Name})
		err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		logger.Info("update-service", lager.Data{"enable-access": service.Name})
		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		return &UpdateServiceOK{Payload: params.Service}

	})

	api.GetDriverInstanceHandler = GetDriverInstanceHandlerFunc(func(params GetDriverInstanceParams, principal interface{}) middleware.Responder {

		instance, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return &GetDriverInstanceNotFound{}
		}
		var conf map[string]interface{}
		err = json.Unmarshal(*instance.Configuration, &conf)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}
		var dials = make([]string, 0)
		for _, dial := range instance.Dials {
			dials = append(dials, dial.ID)
		}

		driverInstance := &genmodel.DriverInstance{
			Configuration: conf,
			Dials:         dials,
			ID:            instance.ID,
			Name:          instance.Name,
			Service:       instance.Service.ID,
		}
		return &GetDriverInstanceOK{Payload: driverInstance}
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams, principal interface{}) middleware.Responder {

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &CreateServicePlanInternalServerError{Payload: err.Error()}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.ID == params.Plan.DialID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return &CreateServicePlanInternalServerError{Payload: err.Error()}
						}

						var plan brokerapi.ServicePlan
						var meta brokerapi.ServicePlanMetadata

						plan.Description = params.Plan.Description
						plan.ID = uuid.NewV4().String()
						plan.Name = params.Plan.Name

						meta.DisplayName = params.Plan.Name
						plan.Metadata = &meta

						dial.Plan = plan
						err = configProvider.SetDial(instance.ID, dial)
						if err != nil {
							return &CreateServicePlanInternalServerError{Payload: err.Error()}
						}
						params.Plan.ID = plan.ID
						return &CreateServicePlanCreated{Payload: params.Plan}
					}
				}
			}
		}

		return &CreateServicePlanInternalServerError{Payload: fmt.Sprintf("Dial %s not found", params.Plan.DialID)}
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation updateCatalog has not yet been implemented")
	})

}
