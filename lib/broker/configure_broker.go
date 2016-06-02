package broker

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/broker/operations"
	"github.com/hpcloud/cf-usb/lib/broker/operations/catalog"
	"github.com/hpcloud/cf-usb/lib/broker/operations/service_instances"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager"
)

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
		return nil, errors.NotImplemented("basic auth  (Basic) has not yet been implemented")
	}

	api.GetServiceInstancesInstanceIDLastOperationHandler = operations.GetServiceInstancesInstanceIDLastOperationHandlerFunc(func(params operations.GetServiceInstancesInstanceIDLastOperationParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation .GetServiceInstancesInstanceIDLastOperation has not yet been implemented")
	})
	api.CatalogCatalogHandler = catalog.CatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation catalog.Catalog has not yet been implemented")
	})
	api.ServiceInstancesCreateServiceInstanceHandler = service_instances.CreateServiceInstanceHandlerFunc(func(params service_instances.CreateServiceInstanceParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.CreateServiceInstance has not yet been implemented")
	})
	api.ServiceInstancesDeprovisionServiceInstanceHandler = service_instances.DeprovisionServiceInstanceHandlerFunc(func(params service_instances.DeprovisionServiceInstanceParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.DeprovisionServiceInstance has not yet been implemented")
	})
	api.ServiceInstancesServiceBindHandler = service_instances.ServiceBindHandlerFunc(func(params service_instances.ServiceBindParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.ServiceBind has not yet been implemented")
	})
	api.ServiceInstancesServiceUnbindHandler = service_instances.ServiceUnbindHandlerFunc(func(params service_instances.ServiceUnbindParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.ServiceUnbind has not yet been implemented")
	})
	api.ServiceInstancesUpdateServiceInstanceHandler = service_instances.UpdateServiceInstanceHandlerFunc(func(params service_instances.UpdateServiceInstanceParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation service_instances.UpdateServiceInstance has not yet been implemented")
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
