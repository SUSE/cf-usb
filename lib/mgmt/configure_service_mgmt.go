package mgmt

import (
	goerrors "errors"
	"github.com/fatih/structs"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
)

func ConfigureServiceAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger) {
	log := logger.Session("usb-mgmt-service")

	api.GetServiceByInstanceIDHandler = GetServiceByInstanceIDHandlerFunc(func(params GetServiceByInstanceIDParams, principal interface{}) middleware.Responder {
		log := log.Session("get-services")
		log.Info("request", lager.Data{"instance-id": params.InstanceID})

		di, _, err := configProvider.GetInstance(params.InstanceID)
		if err != nil {
			return &GetServiceByInstanceIDInternalServerError{Payload: err.Error()}
		}

		service := &genmodel.Service{
			ID:          di.Service.ID,
			InstanceID:  &params.InstanceID,
			Bindable:    di.Service.Bindable,
			Description: di.Service.Description,
			Name:        &di.Service.Name,
			Tags:        di.Service.Tags,
		}

		if di.Service.Metadata != nil {
			service.Metadata = structs.Map(*di.Service.Metadata)
		}

		return &GetServiceByInstanceIDOK{Payload: service}
	})

	api.GetServiceHandler = GetServiceHandlerFunc(func(params GetServiceParams, principal interface{}) middleware.Responder {
		log := log.Session("get-service")
		log.Info("request", lager.Data{"service-id": params.ServiceID})

		serviceInfo, instanceID, err := configProvider.GetService(params.ServiceID)
		if err != nil {
			return &GetServiceInternalServerError{Payload: err.Error()}
		}
		if serviceInfo == nil {
			return &GetServiceNotFound{}
		}

		svc := &genmodel.Service{
			Bindable:    serviceInfo.Bindable,
			InstanceID:  &instanceID,
			ID:          serviceInfo.ID,
			Name:        &serviceInfo.Name,
			Description: serviceInfo.Description,
			Tags:        serviceInfo.Tags,
		}

		if serviceInfo.Metadata != nil {
			svc.Metadata = structs.Map(*serviceInfo.Metadata)
		}

		return &GetServiceOK{Payload: svc}
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams, principal interface{}) middleware.Responder {
		log := log.Session("update-service")
		log.Info("request", lager.Data{"service-id": params.ServiceID})

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		service, instanceid, err := configProvider.GetService(params.ServiceID)
		if err != nil {
			return &UpdateServiceInternalServerError{Payload: err.Error()}
		}

		if service == nil {
			return &UpdateServiceNotFound{}
		}

		if service.Name != *params.Service.Name && *params.Service.Name != "" {
			exists := ccServiceBroker.CheckServiceNameExists(*params.Service.Name)

			if exists == true {
				err := goerrors.New("Service update name parameter validation failed - duplicate naming eror")
				log.Error("update-service-name-validation", err, lager.Data{"Name validation failed for name": params.Service.Name})
				return &UpdateServiceConflict{}
			}

			service.Name = *params.Service.Name
		}

		service.Bindable = params.Service.Bindable
		service.Description = params.Service.Description
		if len(params.Service.Tags) > 0 {
			service.Tags = params.Service.Tags
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

}
