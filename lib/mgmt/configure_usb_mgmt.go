package mgmt

import (
	"encoding/json"
	
	"github.com/fatih/structs"
	
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	. "github.com/hpcloud/cf-usb/lib/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func ConfigureAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface, config *config.Config) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams) (*genmodel.Driver, error) {
		return nil, errors.NotImplemented("operation updateDriver has not yet been implemented")
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams) error {
		return errors.NotImplemented("operation uploadDriver has not yet been implemented")
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func(params DeleteServicePlanParams) error {
		return errors.NotImplemented("operation deleteServicePlan has not yet been implemented")
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams) (*genmodel.DriverSchema, error) {
		return nil, errors.NotImplemented("operation getDriverSchema has not yet been implemented")
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams) (*genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}
		
		for _, d := range config.Drivers{
				if d.ID == params.DriverID{
					
					var instances = make([]string, 0)
					for _, instance := range d.DriverInstances{
						instances = append(instances, instance.ID)
					}					
					
					driver := &genmodel.Driver{
						ID: d.ID,
						DriverType: d.DriverType,	
						DriverInstances: instances,					
					}
					
					return driver, nil
				}
		}
		
		return nil, errors.NotFound("Driver ID: %s not found", params.DriverID)
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(params GetDriversParams) (*[]genmodel.Driver, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}
		
		var drivers = make([]genmodel.Driver, 0)
		
		for _, d := range config.Drivers{
			var instances = make([]string, 0)
			for _, instance := range d.DriverInstances{
				instances = append(instances, instance.ID)
			}					
			
			driver := genmodel.Driver{
				ID: d.ID,
				DriverType: d.DriverType,	
				DriverInstances: instances,					
			}
			
			drivers = append(drivers, driver)
		}
		
		return &drivers, nil
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams) (*genmodel.Driver, error) {
		return nil, errors.NotImplemented("operation createDriver has not yet been implemented")
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams) (*genmodel.Service, error) {
		return nil, errors.NotImplemented("operation updateService has not yet been implemented")
	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams) (*[]genmodel.DriverInstance, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}
		
		var driverInstances = make([]genmodel.DriverInstance, 0)
		
		for _, d := range config.Drivers{
			if params.DriverID != "" && d.ID != params.DriverID{
				continue
			}
			
			for _, di := range d.DriverInstances{
				var dials = make([]string, 0)
				for _, dial := range di.Dials{
					dials = append(dials, dial.ID)
				}
				
				var conf map[string]interface{}
				err := json.Unmarshal(*di.Configuration, &conf)
				if err != nil {
					return nil, err
				}
				
				driverInstance := genmodel.DriverInstance{
					Configuration: conf,
					Dials: dials,
					ID: di.ID,
					DriverID: d.ID,
					Name: di.Name,
					Service: di.Service.ID,
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
		
		for _, d := range config.Drivers{
			for _, di := range d.DriverInstances{
				for _, dial := range di.Dials{
					if dial.Plan.ID == params.PlanID{
						plan := genmodel.Plan{
							Name: dial.Plan.Name,
							ID: dial.Plan.ID,
							DialID: dial.ID,
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
		return nil, errors.NotImplemented("operation createDial has not yet been implemented")
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams) (*genmodel.Dial, error) {
		return nil, errors.NotImplemented("operation updateDial has not yet been implemented")
	})

	api.DeleteServiceHandler = DeleteServiceHandlerFunc(func(params DeleteServiceParams) error {
		return errors.NotImplemented("operation deleteService has not yet been implemented")
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
		
		for _, d := range config.Drivers{			
			for _, di := range d.DriverInstances{
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID{
					continue
				}
				
				service := genmodel.Service{
					ID: di.Service.ID,
					DriverInstanceID: di.ID,
					Bindable: di.Service.Bindable,
					Name: di.Service.Name,
					Tags: di.Service.Tags,
					Metadata: structs.Map(di.Service.Metadata),
				}
				
				services = append(services, service)
			}
		}
		
		return &services, nil
	})

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams) (*genmodel.DriverInstance, error) {
		return nil, errors.NotImplemented("operation createDriverInstance has not yet been implemented")
	})

	api.UpdateDriverInstanceHandler = UpdateDriverInstanceHandlerFunc(func(params UpdateDriverInstanceParams) (*genmodel.DriverInstance, error) {
		return nil, errors.NotImplemented("operation updateDriverInstance has not yet been implemented")
	})

	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams) error {
		return errors.NotImplemented("operation deleteDriver has not yet been implemented")
	})

	api.DeleteDriverInstanceHandler = DeleteDriverInstanceHandlerFunc(func(params DeleteDriverInstanceParams) error {
		return errors.NotImplemented("operation deleteDriverInstance has not yet been implemented")
	})

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams) (*[]genmodel.Dial, error) {
		err := auth.IsAuthenticated(params.Authorization)
		if err != nil {
			return nil, err
		}
		
		var dials = make([]genmodel.Dial, 0)
		
		for _, d := range config.Drivers{			
			for _, di := range d.DriverInstances{
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID{
					continue
				}
				
				for _, dia := range di.Dials{
					
					var conf map[string]interface{}
					err := json.Unmarshal(*dia.Configuration, &conf)
					if err != nil {
						return nil, err
					}
					
					dial := genmodel.Dial{
						Configuration: conf,
						DriverInstanceID: di.ID,
						ID: dia.ID,
						Plan: dia.Plan.ID,
					}
					
					dials = append(dials, dial)
				}
			}
		}
		
		return &dials, nil
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams) (*genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation updateServicePlan has not yet been implemented")
	})

	api.GetBrokerInfoHandler = GetBrokerInfoHandlerFunc(func(params GetBrokerInfoParams) (*genmodel.Broker, error) {

		err := auth.IsAuthenticated(params.Authorization)
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
		
		for _, driver := range config.Drivers{
			for _, instance := range driver.DriverInstances{
				if instance.Service.ID == params.ServiceID{
					svc := &genmodel.Service{
						Bindable: instance.Service.Bindable,
						DriverInstanceID: instance.ID,
						ID: instance.Service.ID,
						Name: instance.Service.Name,
						Description: instance.Service.Description,
						Tags: instance.Service.Tags,
						Metadata: structs.Map(instance.Service.Metadata),
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

		config.BrokerAPI.Credentials.Password = params.Broker.Credentials.Password
		config.BrokerAPI.Credentials.Username = params.Broker.Credentials.Username
		config.BrokerAPI.Listen = params.Broker.Listen

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
		return errors.NotImplemented("operation deleteDial has not yet been implemented")
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func(params GetInfoParams) (*genmodel.Info, error) {
		err := auth.IsAuthenticated(params.Authorization)
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
		
		for _, d := range config.Drivers{			
			for _, di := range d.DriverInstances{		
				if params.DriverInstanceID != "" && di.ID != params.DriverInstanceID{
					continue
				}		
				
				for _, dia := range di.Dials{
					plan := genmodel.Plan{
						Name: dia.Plan.Name,
						ID: dia.Plan.ID,
						DialID: dia.ID,
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
		
		for _, d := range config.Drivers{
			for _, di := range d.DriverInstances{
				for _, dia := range di.Dials{
					if dia.ID == params.DialID{
						
						var conf map[string]interface{}
						err := json.Unmarshal(*dia.Configuration, &conf)
						if err != nil {
							return nil, err
						}
						
						dial := genmodel.Dial{
							Configuration: conf,
							DriverInstanceID: di.ID,
							ID: dia.ID,
							Plan: dia.Plan.ID,
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
		
		for _, d := range config.Drivers{
			for _, di := range d.DriverInstances{
				if di.ID == params.DriverInstanceID{
					var dials = make([]string, 0)
					for _, dial := range di.Dials{
						dials = append(dials, dial.ID)
					}
			
					var conf map[string]interface{}
					err := json.Unmarshal(*di.Configuration, &conf)
					if err != nil {
						return nil, err
					}
					
					driverInstance := genmodel.DriverInstance{
						Configuration: conf,
						Dials: dials,
						ID: di.ID,
						DriverID: d.ID,
						Name: di.Name,
						Service: di.Service.ID,
					}
					
					return &driverInstance, nil
				}
			}
		}
		
		return nil, errors.NotFound("Driver Instance with ID: %s not found", params.DriverInstanceID)
	})

	api.CreateServiceHandler = CreateServiceHandlerFunc(func(params CreateServiceParams) (*genmodel.Service, error) {
		return nil, errors.NotImplemented("operation createService has not yet been implemented")
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams) (*genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation createServicePlan has not yet been implemented")
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(params UpdateCatalogParams) error {
		return errors.NotImplemented("operation updateCatalog has not yet been implemented")
	})

}
