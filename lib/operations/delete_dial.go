package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// DeleteDialHandlerFunc turns a function with the right signature into a delete dial handler
type DeleteDialHandlerFunc func(DeleteDialParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteDialHandlerFunc) Handle(params DeleteDialParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteDialHandler interface for that can handle valid delete dial params
type DeleteDialHandler interface {
	Handle(DeleteDialParams, interface{}) middleware.Responder
}

// NewDeleteDial creates a new http.Handler for the delete dial operation
func NewDeleteDial(ctx *middleware.Context, handler DeleteDialHandler) *DeleteDial {
	return &DeleteDial{Context: ctx, Handler: handler}
}

/*DeleteDial swagger:route DELETE /dials/{dial_id} deleteDial

Delets the `dial` with the **dial_id**

*/
type DeleteDial struct {
	Context *middleware.Context
	Handler DeleteDialHandler
}

func (o *DeleteDial) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewDeleteDialParams()

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
