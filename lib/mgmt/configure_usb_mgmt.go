package mgmt

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/SUSE/cf-usb/lib/brokermodel"
	"github.com/SUSE/cf-usb/lib/config"
	"github.com/SUSE/cf-usb/lib/csm"
	"github.com/SUSE/cf-usb/lib/genmodel"
	"github.com/SUSE/cf-usb/lib/mgmt/authentication"
	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/ccapi"
	"github.com/SUSE/cf-usb/lib/mgmt/operations"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/pivotal-golang/lager"
)

// This file is safe to edit. Once it exists it will not be overwritten

const defaultBrokerName ccapi.BrokerName = "usb"

//ConfigureAPI configures UsbMgmtApi with Interface, config Provider, USBServiceBroker, Logger and a version string
func ConfigureAPI(api *operations.UsbMgmtAPI, auth authentication.Authentication,
	configProvider config.Provider, ccServiceBroker ccapi.USBServiceBroker, csmClient csm.CSM,
	logger lager.Logger, usbVersion string) http.Handler {

	// configure the api here
	log := logger.Session("usb-mgmt")

	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AuthorizationAuth = func(token string) (interface{}, error) {
		err := auth.IsAuthenticated(token)

		if err != nil {
			return nil, err
		}

		log.Debug("authentication-succeeded")

		return true, nil
	}

	api.GetDriverEndpointHandler = operations.GetDriverEndpointHandlerFunc(func(params operations.GetDriverEndpointParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-endpoint")
		log.Info("request", lager.Data{"instance-id": params.DriverEndpointID})

		endpoint, _, err := configProvider.GetInstance(params.DriverEndpointID)
		if err != nil {
			return &operations.GetDriverEndpointsInternalServerError{Payload: err.Error()}
		}
		if endpoint == nil {
			return &operations.GetDriverEndpointNotFound{}
		}

		driverEndpoint := &genmodel.DriverEndpoint{
			ID:                params.DriverEndpointID,
			Name:              &endpoint.Name,
			EndpointURL:       endpoint.TargetURL,
			AuthenticationKey: endpoint.AuthenticationKey,
			Metadata:          map[string]string(endpoint.Service.Metadata),
		}

		return &operations.GetDriverEndpointOK{Payload: driverEndpoint}
	})

	api.GetDriverEndpointsHandler = operations.GetDriverEndpointsHandlerFunc(func(principal interface{}) middleware.Responder {

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &operations.GetDriverEndpointsInternalServerError{Payload: err.Error()}
		}

		var response []*genmodel.DriverEndpoint
		for id, endpoint := range config.Instances {

			var name string
			name = endpoint.Name

			driverEndpoint := &genmodel.DriverEndpoint{
				ID:                id,
				Name:              &name,
				EndpointURL:       endpoint.TargetURL,
				AuthenticationKey: endpoint.AuthenticationKey,
				SkipSSLValidation: &endpoint.SkipSsl,
				CaCertificate:     endpoint.CaCert,
				Metadata:          map[string]string(endpoint.Service.Metadata),
			}

			response = append(response, driverEndpoint)
		}
		return &operations.GetDriverEndpointsOK{Payload: response}
	})

	api.GetInfoHandler = operations.GetInfoHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("get-info")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &operations.GetInfoInternalServerError{
				Payload: err.Error(),
			}
		}

		info := &genmodel.Info{
			BrokerAPIVersion: &config.APIVersion,
			UsbVersion:       &usbVersion,
		}

		return &operations.GetInfoOK{
			Payload: info,
		}
	})

	api.PingDriverEndpointHandler = operations.PingDriverEndpointHandlerFunc(func(params operations.PingDriverEndpointParams, principal interface{}) middleware.Responder {
		// TODO Implement Ping

		return &operations.PingDriverEndpointOK{}
	})

	api.RegisterDriverEndpointHandler = operations.RegisterDriverEndpointHandlerFunc(func(params operations.RegisterDriverEndpointParams, principal interface{}) middleware.Responder {
		log := log.Session("register-driver-endpoint")
		log.Info("request", lager.Data{"id": params.DriverEndpoint.ID, "driver-endpoint-name": params.DriverEndpoint.Name, "driver-endpoint-url": params.DriverEndpoint.EndpointURL})

		if strings.ContainsAny(*params.DriverEndpoint.Name, " ") {
			return &operations.RegisterDriverEndpointInternalServerError{Payload: fmt.Sprintf("Driver endpoint name cannot contain spaces")}
		}

		var instance config.Instance

		instanceID := uuid.NewV4().String()

		params.DriverEndpoint.ID = instanceID

		instance.TargetURL = params.DriverEndpoint.EndpointURL
		instance.AuthenticationKey = params.DriverEndpoint.AuthenticationKey
		if params.DriverEndpoint.SkipSSLValidation == nil {
			instance.SkipSsl = false
		} else {
			instance.SkipSsl = *params.DriverEndpoint.SkipSSLValidation
		}

		instance.CaCert = params.DriverEndpoint.CaCertificate

		driverInstanceNameExist, err := configProvider.InstanceNameExists(*params.DriverEndpoint.Name)
		if err != nil {
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		serviceNameExist, err := ccServiceBroker.CheckServiceNameExists(ccapi.ServiceName(*params.DriverEndpoint.Name))
		if err != nil {
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		var createInstanceError error
		if driverInstanceNameExist && serviceNameExist {
			createInstanceError = fmt.Errorf("a driver instance with the same name already exists")
		} else if driverInstanceNameExist && !serviceNameExist {
			createInstanceError = fmt.Errorf("a driver instance with the same name already exists but is not registered with the Cloud Controller, this service may have been purged - please consider deleting it from the USB")
		} else if !driverInstanceNameExist && serviceNameExist {
			createInstanceError = fmt.Errorf("a service with the same name is already registered with the Cloud Controller")
		}

		if createInstanceError != nil {
			log.Error("check-driver-instance-name-exist", createInstanceError)
			return &operations.RegisterDriverEndpointConflict{Payload: createInstanceError.Error()}
		}

		instance.Name = *params.DriverEndpoint.Name

		err = csmClient.Login(instance.TargetURL, instance.AuthenticationKey, instance.CaCert, instance.SkipSsl)
		if err != nil {
			log.Error("csm-login-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		log.Debug("get-status-information", lager.Data{"url": instance.TargetURL})
		serviceType, err := csmClient.GetStatus()
		if err != nil {
			log.Error("csm-get-status", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		err = configProvider.SetInstance(instanceID, instance)
		if err != nil {
			log.Error("set-driver-instance-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		var defaultDial config.Dial

		defaultDialID := uuid.NewV4().String()

		var plan brokermodel.Plan

		plan.ID = uuid.NewV4().String()
		plan.Description = "default plan"
		plan.Name = "default"
		plan.Free = true

		var meta brokermodel.PlanMetadata

		meta.Name = "default plan"
		meta.Description = "default plan"
		meta.Metadata = struct{ DisplayName string }{"default plan"}

		plan.Metadata = &meta

		defaultDial.Plan = plan
		defaultDialConfig := json.RawMessage([]byte("{}"))
		defaultDial.Configuration = &defaultDialConfig

		err = configProvider.SetDial(instanceID, defaultDialID, defaultDial)
		if err != nil {
			log.Error("set-dial-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		var service brokermodel.CatalogService

		service.ID = uuid.NewV4().String()
		service.Name = *params.DriverEndpoint.Name
		service.Description = "Default service"
		service.Tags = []string{service.Name}
		service.Bindable = true
		service.Metadata = map[string]string(params.DriverEndpoint.Metadata)

		if serviceType == "routing" {
			service.Requires = []string{"route_forwarding"}
		}

		err = configProvider.SetService(instanceID, service)
		if err != nil {
			log.Error("set-service-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			log.Error("load-configuration-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = ccapi.BrokerName(config.ManagementAPI.BrokerName)
		}

		brokerGUID, err := ccServiceBroker.GetServiceBrokerGUIDByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}
		log.Info("create-or-update-service-broker", lager.Data{"guid": brokerGUID})

		if brokerGUID == "" {
			err = ccServiceBroker.Create(brokerName, config.BrokerAPI.ExternalURL, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		} else {
			err = ccServiceBroker.Update(brokerGUID, brokerName, config.BrokerAPI.ExternalURL, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		}
		if err != nil {
			log.Error("create-or-update-service-broker-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		serviceGUID, err := ccServiceBroker.GetServiceGUIDByName(ccapi.ServiceName(*params.DriverEndpoint.Name))
		if err != nil {
			log.Error("get-service-guid-by-name-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		err = ccServiceBroker.EnableServiceAccess(serviceGUID)
		if err != nil {
			log.Error("enable-service-access-failed", err)
			return &operations.RegisterDriverEndpointInternalServerError{Payload: err.Error()}
		}

		return &operations.RegisterDriverEndpointCreated{Payload: params.DriverEndpoint}
	})

	api.UnregisterDriverInstanceHandler = operations.UnregisterDriverInstanceHandlerFunc(func(params operations.UnregisterDriverInstanceParams, principal interface{}) middleware.Responder {
		log := log.Session("unregister-instance")
		log.Info("request", lager.Data{"driver-instance-id": params.DriverEndpointID})

		instance, _, err := configProvider.GetInstance(params.DriverEndpointID)
		if err != nil {
			return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if instance == nil {
			return &operations.UnregisterDriverInstanceNotFound{}
		}
		if ccServiceBroker.CheckServiceInstancesExist(ccapi.ServiceName(instance.Service.Name)) == true {
			return &operations.UnregisterDriverInstanceInternalServerError{Payload: fmt.Sprintf("Cannot delete instance '%s', it still has provisioned service instances", instance.Name)}
		}
		err = configProvider.DeleteInstance(params.DriverEndpointID)
		if err != nil {
			return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
		}

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
		}

		instanceCount := 0
		for _ = range config.Instances {
			instanceCount++
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = ccapi.BrokerName(config.ManagementAPI.BrokerName)
		}

		if instanceCount == 0 {
			err := ccServiceBroker.Delete(brokerName)
			if err != nil {
				log.Error("delete-service-broker-failed", err)
				return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
			}
		} else {
			guid, err := ccServiceBroker.GetServiceBrokerGUIDByName(brokerName)
			if err != nil {
				log.Error("get-service-broker-failed", err)
				return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
			}

			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalURL, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
			if err != nil {
				log.Error("update-service-broker-failed", err)
				return &operations.UnregisterDriverInstanceInternalServerError{Payload: err.Error()}
			}
		}

		return &operations.UnregisterDriverInstanceNoContent{}
	})

	api.UpdateCatalogHandler = operations.UpdateCatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("update-catalog")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &operations.UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = ccapi.BrokerName(config.ManagementAPI.BrokerName)
		}

		brokerGUID, err := ccServiceBroker.GetServiceBrokerGUIDByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &operations.UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		if brokerGUID == "" {
			return &operations.UpdateCatalogInternalServerError{Payload: fmt.Sprintf("Broker %s guid not found", brokerName)}
		}
		err = ccServiceBroker.Update(brokerGUID, brokerName, config.BrokerAPI.ExternalURL, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
		if err != nil {
			return &operations.UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		for _, instance := range config.Instances {
			serviceGUID, err := ccServiceBroker.GetServiceGUIDByName(ccapi.ServiceName(instance.Name))
			if err != nil {
				log.Error("get-service-guid-by-name-failed", err)
				return &operations.UpdateCatalogInternalServerError{Payload: err.Error()}
			}
			err = ccServiceBroker.EnableServiceAccess(serviceGUID)
			if err != nil {
				log.Error("enable-service-access-failed", err)
				return &operations.UpdateCatalogInternalServerError{Payload: err.Error()}
			}
		}

		return &operations.UpdateCatalogOK{}
	})

	api.UpdateDriverEndpointHandler = operations.UpdateDriverEndpointHandlerFunc(func(params operations.UpdateDriverEndpointParams, principal interface{}) middleware.Responder {
		log := log.Session("update-driver-endpoint")
		log.Info("request", lager.Data{"driver-endpoint-id": params.DriverEndpointID})

		instanceInfo, _, err := configProvider.GetInstance(params.DriverEndpointID)
		if err != nil {
			return &operations.UpdateDriverEndpointInternalServerError{Payload: err.Error()}
		}
		if instanceInfo == nil {
			return &operations.UpdateDriverEndpointNotFound{}
		}

		instance := *instanceInfo

		if instanceInfo.Name != *params.DriverEndpoint.Name {
			driverInstanceNameExist, err := configProvider.InstanceNameExists(*params.DriverEndpoint.Name)
			if err != nil {
				return &operations.UpdateDriverEndpointInternalServerError{Payload: err.Error()}
			}

			if driverInstanceNameExist {
				err := fmt.Errorf("A driver instance with the same name already exists")
				log.Error("check-driver-instance-name-exist", err)
				return &operations.UpdateDriverEndpointConflict{}
			}
		}

		if params.DriverEndpoint.AuthenticationKey != "" {
			instance.AuthenticationKey = params.DriverEndpoint.AuthenticationKey
		}
		if params.DriverEndpoint.EndpointURL != "" {
			instance.TargetURL = params.DriverEndpoint.EndpointURL
		}

		instance.Service.Metadata = map[string]string(params.DriverEndpoint.Metadata)

		err = configProvider.SetInstance(params.DriverEndpointID, instance)
		if err != nil {
			return &operations.UpdateDriverEndpointInternalServerError{Payload: err.Error()}
		}

		driverEndpoint := &genmodel.DriverEndpoint{
			ID:                params.DriverEndpointID,
			Name:              &instance.Name,
			EndpointURL:       instance.TargetURL,
			AuthenticationKey: instance.AuthenticationKey,
			Metadata:          map[string]string(params.DriverEndpoint.Metadata),
		}

		return &operations.UpdateDriverEndpointOK{Payload: driverEndpoint}
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
