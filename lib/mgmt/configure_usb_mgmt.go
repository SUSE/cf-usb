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

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams) (*genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var driver config.Driver

		driver.DriverType = params.Driver.DriverType
		driver.ID = params.Driver.ID

		err = configProvider.SetDriver(driver)
		if err != nil {
			return nil, err
		}

		return &params.Driver, nil
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams) error {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return err
		}

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return err
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

		f, err := os.OpenFile(driverPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		defer f.Close()

		reader := bufio.NewReader(params.File.Data)

		sha1 := sha1.New()
		_, err = io.Copy(f, io.TeeReader(reader, sha1))
		if err != nil {
			return err
		}

		sha := base64.StdEncoding.EncodeToString(sha1.Sum(nil))
		if sha != params.Sha {
			f.Close()
			os.Remove(driverPath)
			return errors.New(400, "Checksum mismatch. Expected: %s, got %s", params.Sha, sha)
		}

		return nil
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func(params DeleteServicePlanParams) error {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return err
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.Plan.ID == params.PlanID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return err
						}
						return nil
					}
				}
			}
		}

		return errors.NotFound("Plan %s not found", params.PlanID)
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams) (*genmodel.DriverSchema, error) {
		return nil, errors.NotImplemented("operation getDriverSchema has not yet been implemented")
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams) (*genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		d, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return nil, err
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
		return driver, nil
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(params GetDriversParams) (*[]genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var drivers = make([]genmodel.Driver, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		for _, d := range config.Drivers {
			var instances = make([]string, 0)
			for _, instance := range d.DriverInstances {
				instances = append(instances, instance.ID)
			}

			driver := genmodel.Driver{
				ID:              d.ID,
				DriverType:      d.DriverType,
				DriverInstances: instances,
			}

			drivers = append(drivers, driver)
		}

		return &drivers, nil
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams) (*genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var driver config.Driver

		driver.DriverType = params.Driver.DriverType

		driver.ID = uuid.NewV4().String()

		err = configProvider.SetDriver(driver)
		if err != nil {
			return nil, err
		}
		params.Driver.ID = driver.ID
		return &params.Driver, nil
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams) (*genmodel.Service, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var service brokerapi.Service

		service.Bindable = params.Service.Bindable
		service.Description = params.Service.Description
		service.ID = params.Service.ID
		service.Name = params.Service.Name
		service.Tags = params.Service.Tags

		err = configProvider.SetService(params.Service.DriverInstanceID, service)
		if err != nil {
			return nil, err
		}

		return &params.Service, nil

	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams) (*[]genmodel.DriverInstance, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var driverInstances = make([]genmodel.DriverInstance, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
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
					return nil, err
				}

				driverInstance := genmodel.DriverInstance{
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

		return &driverInstances, nil
	})

	api.PingDriverInstanceHandler = PingDriverInstanceHandlerFunc(func(params PingDriverInstanceParams) error {
		return errors.NotImplemented("operation pingDriverInstance has not yet been implemented")
	})

	api.GetServicePlanHandler = GetServicePlanHandlerFunc(func(params GetServicePlanParams) (*genmodel.Plan, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				for _, dial := range di.Dials {
					if dial.Plan.ID == params.PlanID {
						plan := genmodel.Plan{
							Name:        dial.Plan.Name,
							ID:          dial.Plan.ID,
							DialID:      dial.ID,
							Description: dial.Plan.Description,
							// TODO add free
						}

						return &plan, nil
					}
				}
			}
		}

		return nil, errors.NotFound("Plan with ID: %s not found", params.PlanID)
	})

	api.CreateDialHandler = CreateDialHandlerFunc(func(params CreateDialParams) (*genmodel.Dial, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var dial config.Dial

		dial.ID = uuid.NewV4().String()
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return nil, err
		}
		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration

		err = configProvider.SetDial(params.Dial.DriverInstanceID, dial)
		if err != nil {
			return nil, err
		}

		params.Dial.ID = dial.ID

		return &params.Dial, nil
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams) (*genmodel.Dial, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var dial config.Dial

		dial.ID = params.Dial.ID
		dialconfig, err := json.Marshal(params.Dial.Configuration)
		if err != nil {
			return nil, err
		}

		configuration := json.RawMessage(dialconfig)
		dial.Configuration = &configuration
		dial.Plan = brokerapi.ServicePlan{ID: params.Dial.Plan}
		err = configProvider.SetDial(params.Dial.DriverInstanceID, dial)
		if err != nil {
			return nil, err
		}
		return &params.Dial, nil
	})

	api.DeleteServiceHandler = DeleteServiceHandlerFunc(func(params DeleteServiceParams) error {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return err
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				if instance.Service.ID == params.ServiceID {
					err = configProvider.DeleteService(instance.ID)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
		return nil
	})

	api.GetDialSchemaHandler = GetDialSchemaHandlerFunc(func(params GetDialSchemaParams) (*genmodel.DialSchema, error) {
		return nil, errors.NotImplemented("operation getDialSchema has not yet been implemented")
	})

	api.GetServicesHandler = GetServicesHandlerFunc(func(params GetServicesParams) (*[]genmodel.Service, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var services = make([]genmodel.Service, 0)

		di, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return nil, err
		}

		service := genmodel.Service{
			ID:               di.Service.ID,
			DriverInstanceID: di.ID,
			Bindable:         di.Service.Bindable,
			Name:             di.Service.Name,
			Tags:             di.Service.Tags,
			Metadata:         structs.Map(di.Service.Metadata),
		}

		services = append(services, service)

		return &services, nil
	})

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams) (*genmodel.DriverInstance, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		existingDriver, err := configProvider.GetDriver(params.DriverInstance.DriverID)
		fmt.Println(existingDriver)
		if err != nil {
			return nil, err
		}

		var instance config.DriverInstance
		instance.ID = uuid.NewV4().String()
		instanceConfig, err := json.Marshal(params.DriverInstance.Configuration)
		if err != nil {
			return nil, err
		}
		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.Name = params.DriverInstance.Name

		driverProvider := lib.NewDriverProvider(existingDriver.DriverType, &instance, logger)
		err = driverProvider.Validate()
		if err != nil {
			return nil, err
		}

		err = configProvider.SetDriverInstance(params.DriverInstance.DriverID, instance)
		if err != nil {
			return nil, err
		}

		var defaultDial config.Dial

		defaultDial.ID = uuid.NewV4().String()

		var plan brokerapi.ServicePlan
		plan.Description = "default plan"
		plan.Name = params.DriverInstance.Name + "_default"

		var meta brokerapi.ServicePlanMetadata
		meta.DisplayName = "default plan"

		plan.Metadata = meta

		defaultDial.Plan = plan
		defaultDialConfig := json.RawMessage([]byte("{}"))
		defaultDial.Configuration = &defaultDialConfig

		err = configProvider.SetDial(instance.ID, defaultDial)
		if err != nil {
			return nil, err
		}

		params.DriverInstance.Dials = append(params.DriverInstance.Dials, defaultDial.ID)
		params.DriverInstance.ID = instance.ID

		var service brokerapi.Service
		service.ID = uuid.NewV4().String()
		service.Name = params.DriverInstance.Name + "-default"
		service.Description = "A default service for driver " + params.DriverInstance.Name
		service.Tags = []string{params.DriverInstance.Name}

		err = configProvider.SetService(instance.ID, service)
		if err != nil {
			return nil, err
		}

		params.DriverInstance.Service = service.ID

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			return nil, err
		}

		if guid == "" {
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			return nil, err
		}

		err = ccServiceBroker.EnableServiceAccess(service.Name)
		if err != nil {
			return nil, err
		}

		return &params.DriverInstance, nil
	})

	api.UpdateDriverInstanceHandler = UpdateDriverInstanceHandlerFunc(func(params UpdateDriverInstanceParams) (*genmodel.DriverInstance, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var instance config.DriverInstance
		instance.ID = params.DriverInstanceID

		instanceConfig, err := json.Marshal(params.DriverConfig.Configuration)
		if err != nil {
			return nil, err
		}
		configuration := json.RawMessage(instanceConfig)
		instance.Configuration = &configuration
		instance.Name = params.DriverConfig.Name

		err = configProvider.SetDriverInstance(params.DriverConfig.DriverID, instance)
		if err != nil {
			return nil, err
		}
		return &params.DriverConfig, nil
	})

	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams) error {
		return configProvider.DeleteDriver(params.DriverID)
	})

	api.DeleteDriverInstanceHandler = DeleteDriverInstanceHandlerFunc(func(params DeleteDriverInstanceParams) error {
		return configProvider.DeleteDriverInstance(params.DriverInstanceID)
	})

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams) (*[]genmodel.Dial, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var dials = make([]genmodel.Dial, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
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
						return nil, err
					}

					dial := genmodel.Dial{
						Configuration:    conf,
						DriverInstanceID: di.ID,
						ID:               dia.ID,
						Plan:             dia.Plan.ID,
					}

					dials = append(dials, dial)
				}
			}
		}

		return &dials, nil
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams) (*genmodel.Plan, error) {

		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
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
							plan.Metadata = meta
							dial.Plan = plan
							err = configProvider.SetDial(instance.ID, dial)

							if err != nil {
								return nil, err
							}
							return &params.Plan, nil
						}
					} else {
						return nil, errors.NotFound("Plan %s not found on dial %s", params.PlanID, params.Plan.DialID)
					}
				}
			}
		}

		return nil, errors.NotFound("Dial %s not found", params.Plan.DialID)
	})

	api.GetBrokerInfoHandler = GetBrokerInfoHandlerFunc(func(params GetBrokerInfoParams) (*genmodel.Broker, error) {

		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		broker := &genmodel.Broker{
			Credentials: struct {
				Password string `json:"password"`
				Username string `json:"username"`
			}{
				Username: config.BrokerAPI.Credentials.Username,
				Password: config.BrokerAPI.Credentials.Password,
			},
			Listen: config.BrokerAPI.Listen,
		}
		return broker, nil
	})

	api.GetServiceHandler = GetServiceHandlerFunc(func(params GetServiceParams) (*genmodel.Service, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
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
						Metadata:         structs.Map(instance.Service.Metadata),
					}

					return svc, nil
				}
			}
		}

		return nil, errors.NotFound("Service ID: %s not found", params.ServiceID)
	})

	api.UpdateBrokerInfoHandler = UpdateBrokerInfoHandlerFunc(func(params UpdateBrokerInfoParams) (*genmodel.Broker, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		config.BrokerAPI.Credentials.Password = params.Broker.Credentials.Password
		config.BrokerAPI.Credentials.Username = params.Broker.Credentials.Username
		config.BrokerAPI.Listen = params.Broker.Listen

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			return nil, err
		}

		if guid == "" {
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			return nil, err
		}

		// 		TODO: restart broker listener

		broker := &genmodel.Broker{
			Credentials: struct {
				Password string `json:"password"`
				Username string `json:"username"`
			}{
				Username: config.BrokerAPI.Credentials.Username,
				Password: config.BrokerAPI.Credentials.Password,
			},
			Listen: config.BrokerAPI.Listen,
		}

		return broker, nil
	})

	api.DeleteDialHandler = DeleteDialHandlerFunc(func(params DeleteDialParams) error {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return err
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.ID == params.DialID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return err
						}
						return nil
					}
				}
			}
		}
		return nil
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func(params GetInfoParams) (*genmodel.Info, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}
		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		info := genmodel.Info{
			Version: config.APIVersion,
		}
		return &info, nil
	})

	api.GetServicePlansHandler = GetServicePlansHandlerFunc(func(params GetServicePlansParams) (*[]genmodel.Plan, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var servicePlans = make([]genmodel.Plan, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID {
					continue
				}

				for _, dia := range di.Dials {
					plan := genmodel.Plan{
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

		return &servicePlans, nil
	})

	api.GetDialHandler = GetDialHandlerFunc(func(params GetDialParams) (*genmodel.Dial, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		for _, d := range config.Drivers {
			for _, di := range d.DriverInstances {
				for _, dia := range di.Dials {
					if dia.ID == params.DialID {

						var conf map[string]interface{}
						err := json.Unmarshal(*dia.Configuration, &conf)
						if err != nil {
							return nil, err
						}

						dial := genmodel.Dial{
							Configuration:    conf,
							DriverInstanceID: di.ID,
							ID:               dia.ID,
							Plan:             dia.Plan.ID,
						}

						return &dial, nil
					}
				}
			}
		}

		return nil, errors.NotFound("Dial with ID: %s not found", params.DialID)
	})

	api.GetDriverInstanceHandler = GetDriverInstanceHandlerFunc(func(params GetDriverInstanceParams) (*genmodel.DriverInstance, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		instance, err := configProvider.GetDriverInstance(params.DriverInstanceID)
		if err != nil {
			return nil, err
		}
		var conf map[string]interface{}
		err = json.Unmarshal(*instance.Configuration, &conf)
		if err != nil {
			return nil, err
		}
		var dials = make([]string, 0)
		for _, dial := range instance.Dials {
			dials = append(dials, dial.ID)
		}

		driverInstance := genmodel.DriverInstance{
			Configuration: conf,
			Dials:         dials,
			ID:            instance.ID,
			Name:          instance.Name,
			Service:       instance.Service.ID,
		}
		return &driverInstance, nil
	})

	api.CreateServiceHandler = CreateServiceHandlerFunc(func(params CreateServiceParams) (*genmodel.Service, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		var service brokerapi.Service

		service.Bindable = params.Service.Bindable
		service.Description = params.Service.Description
		service.ID = uuid.NewV4().String()
		service.Name = params.Service.Name
		service.Tags = params.Service.Tags

		err = configProvider.SetService(params.Service.DriverInstanceID, service)
		if err != nil {
			return nil, err
		}
		params.Service.ID = service.ID
		return &params.Service, nil
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams) (*genmodel.Plan, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				for _, dial := range instance.Dials {
					if dial.ID == params.Plan.DialID {
						err = configProvider.DeleteDial(instance.ID, dial.ID)
						if err != nil {
							return nil, err
						}

						var plan brokerapi.ServicePlan
						var meta brokerapi.ServicePlanMetadata

						plan.Description = params.Plan.Description
						plan.ID = uuid.NewV4().String()
						plan.Name = params.Plan.Name

						meta.DisplayName = params.Plan.Name
						plan.Metadata = meta

						dial.Plan = plan
						err = configProvider.SetDial(instance.ID, dial)
						if err != nil {
							return nil, err
						}
						params.Plan.ID = plan.ID
						return &params.Plan, nil
					}
				}
			}
		}

		return nil, errors.NotFound("Dial %s not found", params.Plan.DialID)
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(params UpdateCatalogParams) error {
		return errors.NotImplemented("operation updateCatalog has not yet been implemented")
	})

}
