package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/httpkit/security"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
)

// NewUsbMgmtAPI creates a new UsbMgmt instance
func NewUsbMgmtAPI(spec *spec.Document) *UsbMgmtAPI {
	o := &UsbMgmtAPI{
		spec:            spec,
		handlers:        make(map[string]map[string]http.Handler),
		formats:         strfmt.Default,
		defaultConsumes: "application/json",
		defaultProduces: "application/json",
		ServerShutdown:  func() {},
	}

	return o
}

/*UsbMgmtAPI Universal service broker management API */
type UsbMgmtAPI struct {
	spec            *spec.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	defaultConsumes string
	defaultProduces string
	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer httpkit.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer httpkit.Producer

	// AuthorizationAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key Authorization provided in the header
	AuthorizationAuth func(string) (interface{}, error)

	// CreateDialHandler sets the operation handler for the create dial operation
	CreateDialHandler CreateDialHandler
	// CreateDriverHandler sets the operation handler for the create driver operation
	CreateDriverHandler CreateDriverHandler
	// CreateDriverInstanceHandler sets the operation handler for the create driver instance operation
	CreateDriverInstanceHandler CreateDriverInstanceHandler
	// DeleteDialHandler sets the operation handler for the delete dial operation
	DeleteDialHandler DeleteDialHandler
	// DeleteDriverHandler sets the operation handler for the delete driver operation
	DeleteDriverHandler DeleteDriverHandler
	// DeleteDriverInstanceHandler sets the operation handler for the delete driver instance operation
	DeleteDriverInstanceHandler DeleteDriverInstanceHandler
	// GetAllDialsHandler sets the operation handler for the get all dials operation
	GetAllDialsHandler GetAllDialsHandler
	// GetDialHandler sets the operation handler for the get dial operation
	GetDialHandler GetDialHandler
	// GetDialSchemaHandler sets the operation handler for the get dial schema operation
	GetDialSchemaHandler GetDialSchemaHandler
	// GetDriverHandler sets the operation handler for the get driver operation
	GetDriverHandler GetDriverHandler
	// GetDriverInstanceHandler sets the operation handler for the get driver instance operation
	GetDriverInstanceHandler GetDriverInstanceHandler
	// GetDriverInstancesHandler sets the operation handler for the get driver instances operation
	GetDriverInstancesHandler GetDriverInstancesHandler
	// GetDriverSchemaHandler sets the operation handler for the get driver schema operation
	GetDriverSchemaHandler GetDriverSchemaHandler
	// GetDriversHandler sets the operation handler for the get drivers operation
	GetDriversHandler GetDriversHandler
	// GetInfoHandler sets the operation handler for the get info operation
	GetInfoHandler GetInfoHandler
	// GetServiceHandler sets the operation handler for the get service operation
	GetServiceHandler GetServiceHandler
	// GetServiceByInstanceIDHandler sets the operation handler for the get service by instance id operation
	GetServiceByInstanceIDHandler GetServiceByInstanceIDHandler
	// GetServicePlanHandler sets the operation handler for the get service plan operation
	GetServicePlanHandler GetServicePlanHandler
	// GetServicePlansHandler sets the operation handler for the get service plans operation
	GetServicePlansHandler GetServicePlansHandler
	// PingDriverInstanceHandler sets the operation handler for the ping driver instance operation
	PingDriverInstanceHandler PingDriverInstanceHandler
	// UpdateCatalogHandler sets the operation handler for the update catalog operation
	UpdateCatalogHandler UpdateCatalogHandler
	// UpdateDialHandler sets the operation handler for the update dial operation
	UpdateDialHandler UpdateDialHandler
	// UpdateDriverHandler sets the operation handler for the update driver operation
	UpdateDriverHandler UpdateDriverHandler
	// UpdateDriverInstanceHandler sets the operation handler for the update driver instance operation
	UpdateDriverInstanceHandler UpdateDriverInstanceHandler
	// UpdateServiceHandler sets the operation handler for the update service operation
	UpdateServiceHandler UpdateServiceHandler
	// UpdateServicePlanHandler sets the operation handler for the update service plan operation
	UpdateServicePlanHandler UpdateServicePlanHandler
	// UploadDriverHandler sets the operation handler for the upload driver operation
	UploadDriverHandler UploadDriverHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()
}

// SetDefaultProduces sets the default produces media type
func (o *UsbMgmtAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *UsbMgmtAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// DefaultProduces returns the default produces media type
func (o *UsbMgmtAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *UsbMgmtAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *UsbMgmtAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *UsbMgmtAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the UsbMgmtAPI
func (o *UsbMgmtAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.AuthorizationAuth == nil {
		unregistered = append(unregistered, "AuthorizationAuth")
	}

	if o.CreateDialHandler == nil {
		unregistered = append(unregistered, "CreateDialHandler")
	}

	if o.CreateDriverHandler == nil {
		unregistered = append(unregistered, "CreateDriverHandler")
	}

	if o.CreateDriverInstanceHandler == nil {
		unregistered = append(unregistered, "CreateDriverInstanceHandler")
	}

	if o.DeleteDialHandler == nil {
		unregistered = append(unregistered, "DeleteDialHandler")
	}

	if o.DeleteDriverHandler == nil {
		unregistered = append(unregistered, "DeleteDriverHandler")
	}

	if o.DeleteDriverInstanceHandler == nil {
		unregistered = append(unregistered, "DeleteDriverInstanceHandler")
	}

	if o.GetAllDialsHandler == nil {
		unregistered = append(unregistered, "GetAllDialsHandler")
	}

	if o.GetDialHandler == nil {
		unregistered = append(unregistered, "GetDialHandler")
	}

	if o.GetDialSchemaHandler == nil {
		unregistered = append(unregistered, "GetDialSchemaHandler")
	}

	if o.GetDriverHandler == nil {
		unregistered = append(unregistered, "GetDriverHandler")
	}

	if o.GetDriverInstanceHandler == nil {
		unregistered = append(unregistered, "GetDriverInstanceHandler")
	}

	if o.GetDriverInstancesHandler == nil {
		unregistered = append(unregistered, "GetDriverInstancesHandler")
	}

	if o.GetDriverSchemaHandler == nil {
		unregistered = append(unregistered, "GetDriverSchemaHandler")
	}

	if o.GetDriversHandler == nil {
		unregistered = append(unregistered, "GetDriversHandler")
	}

	if o.GetInfoHandler == nil {
		unregistered = append(unregistered, "GetInfoHandler")
	}

	if o.GetServiceHandler == nil {
		unregistered = append(unregistered, "GetServiceHandler")
	}

	if o.GetServiceByInstanceIDHandler == nil {
		unregistered = append(unregistered, "GetServiceByInstanceIDHandler")
	}

	if o.GetServicePlanHandler == nil {
		unregistered = append(unregistered, "GetServicePlanHandler")
	}

	if o.GetServicePlansHandler == nil {
		unregistered = append(unregistered, "GetServicePlansHandler")
	}

	if o.PingDriverInstanceHandler == nil {
		unregistered = append(unregistered, "PingDriverInstanceHandler")
	}

	if o.UpdateCatalogHandler == nil {
		unregistered = append(unregistered, "UpdateCatalogHandler")
	}

	if o.UpdateDialHandler == nil {
		unregistered = append(unregistered, "UpdateDialHandler")
	}

	if o.UpdateDriverHandler == nil {
		unregistered = append(unregistered, "UpdateDriverHandler")
	}

	if o.UpdateDriverInstanceHandler == nil {
		unregistered = append(unregistered, "UpdateDriverInstanceHandler")
	}

	if o.UpdateServiceHandler == nil {
		unregistered = append(unregistered, "UpdateServiceHandler")
	}

	if o.UpdateServicePlanHandler == nil {
		unregistered = append(unregistered, "UpdateServicePlanHandler")
	}

	if o.UploadDriverHandler == nil {
		unregistered = append(unregistered, "UploadDriverHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *UsbMgmtAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *UsbMgmtAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]httpkit.Authenticator {

	result := make(map[string]httpkit.Authenticator)
	for name, scheme := range schemes {
		switch name {

		case "Authorization":

			result[name] = security.APIKeyAuth(scheme.Name, scheme.In, func(tok string) (interface{}, error) { return o.AuthorizationAuth(tok) })

		}
	}
	return result

}

// ConsumersFor gets the consumers for the specified media types
func (o *UsbMgmtAPI) ConsumersFor(mediaTypes []string) map[string]httpkit.Consumer {

	result := make(map[string]httpkit.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *UsbMgmtAPI) ProducersFor(mediaTypes []string) map[string]httpkit.Producer {

	result := make(map[string]httpkit.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *UsbMgmtAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

func (o *UsbMgmtAPI) initHandlerCache() {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers[strings.ToUpper("POST")] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/dials"] = NewCreateDial(o.context, o.CreateDialHandler)

	if o.handlers["POST"] == nil {
		o.handlers[strings.ToUpper("POST")] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/drivers"] = NewCreateDriver(o.context, o.CreateDriverHandler)

	if o.handlers["POST"] == nil {
		o.handlers[strings.ToUpper("POST")] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/driver_instances"] = NewCreateDriverInstance(o.context, o.CreateDriverInstanceHandler)

	if o.handlers["DELETE"] == nil {
		o.handlers[strings.ToUpper("DELETE")] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/dials/{dial_id}"] = NewDeleteDial(o.context, o.DeleteDialHandler)

	if o.handlers["DELETE"] == nil {
		o.handlers[strings.ToUpper("DELETE")] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/drivers/{driver_id}"] = NewDeleteDriver(o.context, o.DeleteDriverHandler)

	if o.handlers["DELETE"] == nil {
		o.handlers[strings.ToUpper("DELETE")] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/driver_instances/{driver_instance_id}"] = NewDeleteDriverInstance(o.context, o.DeleteDriverInstanceHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/dials"] = NewGetAllDials(o.context, o.GetAllDialsHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/dials/{dial_id}"] = NewGetDial(o.context, o.GetDialHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/drivers/{driver_id}/dial_schema"] = NewGetDialSchema(o.context, o.GetDialSchemaHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/drivers/{driver_id}"] = NewGetDriver(o.context, o.GetDriverHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/driver_instances/{driver_instance_id}"] = NewGetDriverInstance(o.context, o.GetDriverInstanceHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/driver_instances"] = NewGetDriverInstances(o.context, o.GetDriverInstancesHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/drivers/{driver_id}/config_schema"] = NewGetDriverSchema(o.context, o.GetDriverSchemaHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/drivers"] = NewGetDrivers(o.context, o.GetDriversHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/info"] = NewGetInfo(o.context, o.GetInfoHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/services/{service_id}"] = NewGetService(o.context, o.GetServiceHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/services"] = NewGetServiceByInstanceID(o.context, o.GetServiceByInstanceIDHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/plans/{plan_id}"] = NewGetServicePlan(o.context, o.GetServicePlanHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/plans"] = NewGetServicePlans(o.context, o.GetServicePlansHandler)

	if o.handlers["GET"] == nil {
		o.handlers[strings.ToUpper("GET")] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/driver_instances/{driver_instance_id}/ping"] = NewPingDriverInstance(o.context, o.PingDriverInstanceHandler)

	if o.handlers["POST"] == nil {
		o.handlers[strings.ToUpper("POST")] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/update_catalog"] = NewUpdateCatalog(o.context, o.UpdateCatalogHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/dials/{dial_id}"] = NewUpdateDial(o.context, o.UpdateDialHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/drivers/{driver_id}"] = NewUpdateDriver(o.context, o.UpdateDriverHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/driver_instances/{driver_instance_id}"] = NewUpdateDriverInstance(o.context, o.UpdateDriverInstanceHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/services/{service_id}"] = NewUpdateService(o.context, o.UpdateServiceHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/plans/{plan_id}"] = NewUpdateServicePlan(o.context, o.UpdateServicePlanHandler)

	if o.handlers["PUT"] == nil {
		o.handlers[strings.ToUpper("PUT")] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/drivers/{driver_id}/bits"] = NewUploadDriver(o.context, o.UploadDriverHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *UsbMgmtAPI) Serve(builder middleware.Builder) http.Handler {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}

	return o.context.APIHandler(builder)
}
