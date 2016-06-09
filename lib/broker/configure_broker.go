package broker

import (
	"crypto/tls"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/broker/operations"
	"github.com/hpcloud/cf-usb/lib/broker/operations/catalog"
	"github.com/hpcloud/cf-usb/lib/broker/operations/service_instances"
	"github.com/hpcloud/cf-usb/lib/brokermodel"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager"
)

func getBrokerError(s string) *brokermodel.BrokerError {
	msg := s
	brokerError := brokermodel.BrokerError{Message: &msg}
	return &brokerError
}

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.BrokerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func ConfigureAPI(api *operations.BrokerAPI, csm csm.CSM, configProvider config.Provider, logger lager.Logger) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.BasicAuth = func(user string, pass string) (interface{}, error) {
		conf, err := configProvider.LoadConfiguration()
		if err != nil {
			return nil, err
		}
		if conf.BrokerAPI.Credentials.Username == user &&
			conf.BrokerAPI.Credentials.Password == pass {
			return true, nil
		}
		return false, nil
	}

	api.GetServiceInstancesInstanceIDLastOperationHandler = operations.GetServiceInstancesInstanceIDLastOperationHandlerFunc(func(params operations.GetServiceInstancesInstanceIDLastOperationParams, principal interface{}) middleware.Responder {
		//TODO add async
		exists, err := csm.WorkspaceExists(params.InstanceID)
		payload := &brokermodel.LastOperation{}
		if err != nil {
			payload.State = "failed"
			payload.Description = err.Error()
			return operations.NewGetServiceInstancesInstanceIDLastOperationOK().WithPayload(payload)
		}

		payload.Description = fmt.Sprintf("resources exists = %t", exists)
		payload.State = "succeeded"
		return operations.NewGetServiceInstancesInstanceIDLastOperationOK().WithPayload(payload)

	})
	api.CatalogCatalogHandler = catalog.CatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		var cat = brokermodel.CatalogServices{}

		conf, err := configProvider.LoadConfiguration()

		if err != nil {
			return catalog.NewCatalogDefault(500).WithPayload(getBrokerError(err.Error()))
		}

		for _, instance := range conf.Instances {
			catServ := instance.Service
			for _, dial := range instance.Dials {
				dialTemp := dial
				catServ.Plans = append(catServ.Plans, &dialTemp.Plan)
			}
			cat.Services = append(cat.Services, &catServ)

		}
		return catalog.NewCatalogOK().WithPayload(&cat)
		//return middleware.NotImplemented("operation catalog.Catalog has not yet been implemented")
	})
	api.ServiceInstancesCreateServiceInstanceHandler = service_instances.CreateServiceInstanceHandlerFunc(func(params service_instances.CreateServiceInstanceParams, principal interface{}) middleware.Responder {
		servID, err := getServiceAfterLogin(csm, configProvider, params.Service.ServiceID)
		if err != nil {
			return service_instances.NewCreateServiceInstanceDefault(401).WithPayload(getBrokerError(err.Error()))
		}
		if servID == nil {
			return service_instances.NewCreateServiceInstanceDefault(410).WithPayload(getBrokerError(servID.ID + " not found"))
		}
		exists, err := csm.WorkspaceExists(params.InstanceID)
		if err != nil {
			return service_instances.NewCreateServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
		}
		if exists {
			return service_instances.NewCreateServiceInstanceConflict()
		}

		err = csm.CreateWorkspace(params.InstanceID)
		if err != nil {
			return service_instances.NewCreateServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
		}

		logger.Info("provision-instance-request-completed", lager.Data{"instance-id": params.InstanceID})

		return service_instances.NewCreateServiceInstanceCreated().WithPayload(&brokermodel.DashboardURL{})
		//TODO create async
	})
	api.ServiceInstancesDeprovisionServiceInstanceHandler = service_instances.DeprovisionServiceInstanceHandlerFunc(func(params service_instances.DeprovisionServiceInstanceParams, principal interface{}) middleware.Responder {
		servID, err := getServiceAfterLogin(csm, configProvider, params.DeprovisionParameters.ServiceID)
		if err != nil {
			return service_instances.NewDeprovisionServiceInstanceDefault(401).WithPayload(getBrokerError(err.Error()))
		}
		if servID == nil {
			return service_instances.NewDeprovisionServiceInstanceDefault(410).WithPayload(getBrokerError(servID.ID + " not found"))
		}
		exists, err := csm.WorkspaceExists(params.InstanceID)
		if err != nil {
			return service_instances.NewDeprovisionServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
		}
		if !exists {
			return service_instances.NewDeprovisionServiceInstanceGone()
		}
		err = csm.DeleteWorkspace(params.InstanceID)
		if err != nil {
			return service_instances.NewDeprovisionServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
		}
		logger.Info("deprovision-service--instance-request-completed", lager.Data{"instance-id": params.InstanceID, "service-id": params.DeprovisionParameters.ServiceID})

		return service_instances.NewDeprovisionServiceInstanceOK().WithPayload(&brokermodel.Empty{})
	})
	api.ServiceInstancesServiceBindHandler = service_instances.ServiceBindHandlerFunc(func(params service_instances.ServiceBindParams, principal interface{}) middleware.Responder {

		servID, err := getServiceAfterLogin(csm, configProvider, params.Binding.ServiceID)
		if err != nil {
			return service_instances.NewServiceBindDefault(401).WithPayload(getBrokerError(err.Error()))
		}
		//if no service with this ID was found we send Gone HTTP header
		if servID == nil {
			return service_instances.NewServiceBindDefault(410).WithPayload(getBrokerError(servID.ID + " not found"))
		}

		exists, err := csm.ConnectionExists(params.InstanceID, params.BindingID)
		if err != nil {
			return service_instances.NewServiceBindDefault(500).WithPayload(getBrokerError(err.Error()))
		}
		//if it allready exists we send it the OK - 200 HTTP header
		if exists {
			return service_instances.NewServiceBindOK()
		}
		results, err := csm.CreateConnection(params.InstanceID, params.BindingID)
		if err != nil {
			return service_instances.NewServiceBindDefault(500).WithPayload(getBrokerError(err.Error()))
		}

		bindingResponse := brokermodel.BindingResponse{}
		bindingResponse.Credentials = results

		logger.Info("generate-credentials-request-completed", lager.Data{"binding-id": params.BindingID})

		return service_instances.NewServiceBindCreated().WithPayload(&bindingResponse)

	})
	api.ServiceInstancesServiceUnbindHandler = service_instances.ServiceUnbindHandlerFunc(func(params service_instances.ServiceUnbindParams, principal interface{}) middleware.Responder {
		servID, err := getServiceAfterLogin(csm, configProvider, params.UnbindParameters.ServiceID)
		if err != nil {
			return service_instances.NewServiceUnbindDefault(401).WithPayload(getBrokerError(err.Error()))
		}
		//if no service with this ID was found we send Gone HTTP header
		if servID == nil {
			return service_instances.NewServiceUnbindDefault(410).WithPayload(getBrokerError(servID.ID + " not found"))
		}

		exists, err := csm.ConnectionExists(params.InstanceID, params.BindingID)
		if err != nil {
			return service_instances.NewServiceUnbindDefault(500).WithPayload(getBrokerError(err.Error()))
		}
		//if it does not exist we send it the Gone 410 HTTP header
		if !exists {
			return service_instances.NewServiceUnbindGone()
		}
		err = csm.DeleteConnection(params.InstanceID, params.BindingID)
		if err != nil {
			return service_instances.NewServiceUnbindDefault(500).WithPayload(getBrokerError(err.Error()))
		}

		logger.Info("unbind-instance-completed", lager.Data{"binding-id": params.BindingID})
		r := brokermodel.Empty{}
		return service_instances.NewServiceUnbindOK().WithPayload(&r)
	})
	api.ServiceInstancesUpdateServiceInstanceHandler = service_instances.UpdateServiceInstanceHandlerFunc(func(params service_instances.UpdateServiceInstanceParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.UpdateServiceInstance has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func getServiceAfterLogin(csm csm.CSM, configProvider config.Provider, serviceID string) (*brokermodel.CatalogService, error) {
	conf, err := configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}
	for _, driverInstance := range conf.Instances {
		if driverInstance.Service.ID == serviceID {
			if driverInstance.TargetURL != "" {
				err = csm.Login(driverInstance.TargetURL, driverInstance.AuthenticationKey)
				if err != nil {
					return nil, err
				}
				return &driverInstance.Service, nil
			}
		}
	}
	return nil, nil
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
