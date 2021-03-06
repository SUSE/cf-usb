package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// RegisterDriverEndpointHandlerFunc turns a function with the right signature into a register driver endpoint handler
type RegisterDriverEndpointHandlerFunc func(RegisterDriverEndpointParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn RegisterDriverEndpointHandlerFunc) Handle(params RegisterDriverEndpointParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// RegisterDriverEndpointHandler interface for that can handle valid register driver endpoint params
type RegisterDriverEndpointHandler interface {
	Handle(RegisterDriverEndpointParams, interface{}) middleware.Responder
}

// NewRegisterDriverEndpoint creates a new http.Handler for the register driver endpoint operation
func NewRegisterDriverEndpoint(ctx *middleware.Context, handler RegisterDriverEndpointHandler) *RegisterDriverEndpoint {
	return &RegisterDriverEndpoint{Context: ctx, Handler: handler}
}

/*RegisterDriverEndpoint swagger:route POST /driver_endpoints registerDriverEndpoint

Registers a driver endpoint with the USB

*/
type RegisterDriverEndpoint struct {
	Context *middleware.Context
	Handler RegisterDriverEndpointHandler
}

func (o *RegisterDriverEndpoint) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewRegisterDriverEndpointParams()

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
