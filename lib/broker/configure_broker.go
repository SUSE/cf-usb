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
	"github.com/hpcloud/cf-usb/lib/brokermodel"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager"
)

var (
	brokerCsm            csm.CSM
	brokerConfigProvider config.Provider
	brokerLogger         lager.Logger
)

const (
	failed   string = "failed"
	succeded string = "succeded"
)

func getBrokerError(s string) *brokermodel.BrokerError {
	msg := s
	brokerError := brokermodel.BrokerError{Message: &msg}
	return &brokerError
}

func idLastOperationHandler(params operations.GetServiceInstancesInstanceIDLastOperationParams, principal interface{}) middleware.Responder {
	//TODO add async
	exists, isNoop, err := brokerCsm.WorkspaceExists(params.InstanceID)

	payload := &brokermodel.LastOperation{}

	if err != nil {
		payload.State = failed
		payload.Description = err.Error()
		brokerLogger.Info("last-operation-error", lager.Data{"error": err.Error()})
		return operations.NewGetServiceInstancesInstanceIDLastOperationOK().WithPayload(payload)
	}

	payload.Description = fmt.Sprintf("resources exists = %t; Operation is noop = %t", exists, isNoop)
	payload.State = succeded
	brokerLogger.Info("last-operation-completed", lager.Data{"instance-id": params.InstanceID})
	return operations.NewGetServiceInstancesInstanceIDLastOperationOK().WithPayload(payload)

}

func catalogHandler(principal interface{}) middleware.Responder {
	var cat = brokermodel.CatalogServices{}

	conf, err := brokerConfigProvider.LoadConfiguration()

	if err != nil {
		brokerLogger.Info("catalog-request-error", lager.Data{"error": err.Error()})
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

	brokerLogger.Info("catalog-request-completed", lager.Data{"catalog-services": cat.Services})

	return catalog.NewCatalogOK().WithPayload(&cat)

}

func createServiceInstanceHandler(params operations.CreateServiceInstanceParams, principal interface{}) middleware.Responder {

	servID, err := getServiceAfterLogin(brokerCsm, brokerConfigProvider, params.Service.ServiceID)

	if err != nil {
		brokerLogger.Info("provision-instance-request-error", lager.Data{"error": err.Error()})
		return operations.NewCreateServiceInstanceDefault(401).WithPayload(getBrokerError(err.Error()))
	}

	if servID == nil {
		brokerLogger.Info("provision-instance-request-not-in-catalog", lager.Data{"service-id": params.Service.ServiceID})
		return operations.NewCreateServiceInstanceDefault(404).WithPayload(getBrokerError(params.Service.ServiceID + " not found"))
	}

	exists, isNoop, err := brokerCsm.WorkspaceExists(params.InstanceID)

	if err != nil {
		brokerLogger.Info("provision-instance-request-error", lager.Data{"error": err.Error()})
		return operations.NewCreateServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	if exists && !isNoop {
		brokerLogger.Info("provision-instance-request-conflict", lager.Data{"instance-id": params.InstanceID, "service-id": params.Service.ServiceID})
		return operations.NewCreateServiceInstanceConflict()
	}

	err = brokerCsm.CreateWorkspace(params.InstanceID)

	if err != nil {
		return operations.NewCreateServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	brokerLogger.Info("provision-instance-request-completed", lager.Data{"instance-id": params.InstanceID, "service-id": params.Service.ServiceID})

	return operations.NewCreateServiceInstanceCreated().WithPayload(&brokermodel.DashboardURL{})
	//TODO create async
}

func deprovisionServiceInstanceHandler(params operations.DeprovisionServiceInstanceParams, principal interface{}) middleware.Responder {

	servID, err := getServiceAfterLogin(brokerCsm, brokerConfigProvider, params.ServiceID)

	if err != nil {
		brokerLogger.Info("deprovision-service-error", lager.Data{"error": err.Error()})
		return operations.NewDeprovisionServiceInstanceDefault(401).WithPayload(getBrokerError(err.Error()))
	}

	if servID == nil {
		brokerLogger.Info("deprovision-service-not-in-catalog", lager.Data{"service-id": params.ServiceID})
		return operations.NewDeprovisionServiceInstanceDefault(404).WithPayload(getBrokerError(params.ServiceID + " not found"))
	}

	exists, isNoop, err := brokerCsm.WorkspaceExists(params.InstanceID)

	if err != nil {
		brokerLogger.Info("deprovision-service-error", lager.Data{"error": err.Error()})
		return operations.NewDeprovisionServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	if !exists && !isNoop {
		brokerLogger.Info("deprovision-service-missing", lager.Data{"instance-id": params.InstanceID})
		return operations.NewDeprovisionServiceInstanceDefault(404).WithPayload(getBrokerError("No bind with this name found"))
	}

	err = brokerCsm.DeleteWorkspace(params.InstanceID)

	if err != nil {
		brokerLogger.Info("deprovision-service-error", lager.Data{"error": err.Error()})
		return operations.NewDeprovisionServiceInstanceDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	brokerLogger.Info("deprovision-service-instance-request-completed", lager.Data{"instance-id": params.InstanceID, "service-id": params.ServiceID})

	return operations.NewDeprovisionServiceInstanceOK().WithPayload(map[string]interface{}{})
}

func serviceBindHandler(params operations.ServiceBindParams, principal interface{}) middleware.Responder {

	servID, err := getServiceAfterLogin(brokerCsm, brokerConfigProvider, params.Binding.ServiceID)

	if err != nil {
		brokerLogger.Info("generate-credentials-service-error", lager.Data{"error": err.Error()})
		return operations.NewServiceBindDefault(401).WithPayload(getBrokerError(err.Error()))
	}

	//if no service with this ID was found we send not found HTTP header
	if servID == nil {
		brokerLogger.Info("generate-credentials-service-not-in-catalog", lager.Data{"service-id": params.Binding.ServiceID})
		return operations.NewServiceBindDefault(404).WithPayload(getBrokerError(params.Binding.ServiceID + " not found"))
	}

	exists, isNoop, err := brokerCsm.ConnectionExists(params.InstanceID, params.BindingID)

	if err != nil {
		brokerLogger.Info("generate-credentials-service-error", lager.Data{"error": err.Error()})
		return operations.NewServiceBindDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	//if it already exists we send it the Conflict - 409 HTTP header
	if exists && !isNoop {
		brokerLogger.Info("generate-credentials-service-allready exists", lager.Data{"instance-id": params.InstanceID, "binding-id": params.BindingID})
		return operations.NewServiceBindConflict().WithPayload(map[string]interface{}{})
	}

	results, err := brokerCsm.CreateConnection(params.InstanceID, params.BindingID)

	if err != nil {
		brokerLogger.Info("generate-credentials-service-error", lager.Data{"error": err.Error()})
		return operations.NewServiceBindDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	bindingResponse := brokermodel.BindingResponse{}
	bindingResponse.Credentials = results

	brokerLogger.Info("generate-credentials-request-completed", lager.Data{"binding-id": params.BindingID, "service-id": params.Binding.ServiceID})

	return operations.NewServiceBindCreated().WithPayload(&bindingResponse)

}

func serviceUnbindHandler(params operations.ServiceUnbindParams, principal interface{}) middleware.Responder {

	servID, err := getServiceAfterLogin(brokerCsm, brokerConfigProvider, params.ServiceID)

	if err != nil {
		brokerLogger.Info("unbind-instance-error", lager.Data{"error": err.Error()})
		return operations.NewServiceUnbindDefault(401).WithPayload(getBrokerError(err.Error()))
	}
	//if no service with this ID was found we send Gone HTTP header
	if servID == nil {
		brokerLogger.Info("unbind-instance-not-in-catalog", lager.Data{"service-id": params.ServiceID})
		return operations.NewServiceUnbindDefault(404).WithPayload(getBrokerError(params.ServiceID + " not found"))
	}

	exists, isNoop, err := brokerCsm.ConnectionExists(params.InstanceID, params.BindingID)

	if err != nil {
		brokerLogger.Info("unbind-instance-error", lager.Data{"error": err.Error()})
		return operations.NewServiceUnbindDefault(500).WithPayload(getBrokerError(err.Error()))
	}
	//if it does not exist we send it the Not found 404 HTTP header
	if !exists && !isNoop {
		brokerLogger.Info("unbind-instance-missing", lager.Data{"instance-id": params.InstanceID, "binding-id": params.BindingID})
		return operations.NewServiceUnbindDefault(404).WithPayload(getBrokerError(fmt.Sprintf("Binding %s not found", params.BindingID)))
	}

	err = brokerCsm.DeleteConnection(params.InstanceID, params.BindingID)

	if err != nil {
		brokerLogger.Info("unbind-instance-error", lager.Data{"error": err.Error()})
		return operations.NewServiceUnbindDefault(500).WithPayload(getBrokerError(err.Error()))
	}

	brokerLogger.Info("unbind-instance-completed", lager.Data{"binding-id": params.BindingID})

	return operations.NewServiceUnbindOK().WithPayload(map[string]interface{}{})
}

func updateServiceInstanceHandler(params operations.UpdateServiceInstanceParams, principal interface{}) middleware.Responder {
	brokerLogger.Info("unpdate-service-instance-not-implemented", lager.Data{"instance-id": params.InstanceID})
	return middleware.NotImplemented("operation operations.UpdateServiceInstance has not yet been implemented")
}

func basicAuth(user string, pass string) (interface{}, error) {
	conf, err := brokerConfigProvider.LoadConfiguration()

	if err != nil {
		return nil, err
	}

	if conf.BrokerAPI.Credentials.Username == user &&
		conf.BrokerAPI.Credentials.Password == pass {
		return true, nil
	}

	brokerLogger.Info("Invalid username/pass", lager.Data{"message": "invalid username/password combination"})

	return nil, fmt.Errorf("Not authorized")
}

//ConfigureAPI is the function that defines what functions will handle the requestss
func ConfigureAPI(api *operations.BrokerAPI, csm csm.CSM, configProvider config.Provider, logger lager.Logger) http.Handler {

	brokerCsm = csm
	brokerLogger = logger
	brokerConfigProvider = configProvider

	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.BasicAuth = basicAuth

	api.GetServiceInstancesInstanceIDLastOperationHandler =
		operations.GetServiceInstancesInstanceIDLastOperationHandlerFunc(idLastOperationHandler)

	api.CatalogCatalogHandler =
		catalog.CatalogHandlerFunc(catalogHandler)

	api.CreateServiceInstanceHandler =
		operations.CreateServiceInstanceHandlerFunc(createServiceInstanceHandler)

	api.DeprovisionServiceInstanceHandler =
		operations.DeprovisionServiceInstanceHandlerFunc(deprovisionServiceInstanceHandler)

	api.ServiceBindHandler =
		operations.ServiceBindHandlerFunc(serviceBindHandler)

	api.ServiceUnbindHandler =
		operations.ServiceUnbindHandlerFunc(serviceUnbindHandler)

	api.UpdateServiceInstanceHandler =
		operations.UpdateServiceInstanceHandlerFunc(updateServiceInstanceHandler)

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
				err = csm.Login(driverInstance.TargetURL, driverInstance.AuthenticationKey, driverInstance.CaCert, driverInstance.SkipSsl)
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
