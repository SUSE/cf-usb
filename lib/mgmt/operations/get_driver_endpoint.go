package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetDriverEndpointHandlerFunc turns a function with the right signature into a get driver endpoint handler
type GetDriverEndpointHandlerFunc func(GetDriverEndpointParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDriverEndpointHandlerFunc) Handle(params GetDriverEndpointParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetDriverEndpointHandler interface for that can handle valid get driver endpoint params
type GetDriverEndpointHandler interface {
	Handle(GetDriverEndpointParams, interface{}) middleware.Responder
}

// NewGetDriverEndpoint creates a new http.Handler for the get driver endpoint operation
func NewGetDriverEndpoint(ctx *middleware.Context, handler GetDriverEndpointHandler) *GetDriverEndpoint {
	return &GetDriverEndpoint{Context: ctx, Handler: handler}
}

/*GetDriverEndpoint swagger:route GET /driver_endpoints/{driver_endpoint_id} getDriverEndpoint

Gets details for a specific driver endpoint


*/
type GetDriverEndpoint struct {
	Context *middleware.Context
	Handler GetDriverEndpointHandler
}

func (o *GetDriverEndpoint) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetDriverEndpointParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}