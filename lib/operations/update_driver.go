package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// UpdateDriverHandlerFunc turns a function with the right signature into a update driver handler
type UpdateDriverHandlerFunc func(UpdateDriverParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateDriverHandlerFunc) Handle(params UpdateDriverParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// UpdateDriverHandler interface for that can handle valid update driver params
type UpdateDriverHandler interface {
	Handle(UpdateDriverParams, interface{}) middleware.Responder
}

// NewUpdateDriver creates a new http.Handler for the update driver operation
func NewUpdateDriver(ctx *middleware.Context, handler UpdateDriverHandler) *UpdateDriver {
	return &UpdateDriver{Context: ctx, Handler: handler}
}

/*UpdateDriver swagger:route PUT /drivers/{driver_id} updateDriver

Update driver

*/
type UpdateDriver struct {
	Context *middleware.Context
	Params  UpdateDriverParams
	Handler UpdateDriverHandler
}

func (o *UpdateDriver) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
