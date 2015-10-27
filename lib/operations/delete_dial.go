package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// DeleteDialHandlerFunc turns a function with the right signature into a delete dial handler
type DeleteDialHandlerFunc func(DeleteDialParams) error

func (fn DeleteDialHandlerFunc) Handle(params DeleteDialParams) error {
	return fn(params)
}

// DeleteDialHandler interface for that can handle valid delete dial params
type DeleteDialHandler interface {
	Handle(DeleteDialParams) error
}

// NewDeleteDial creates a new http.Handler for the delete dial operation
func NewDeleteDial(ctx *middleware.Context, handler DeleteDialHandler) *DeleteDial {
	return &DeleteDial{Context: ctx, Handler: handler}
}

/*
Delets the `dial` with the **dial_id**
*/
type DeleteDial struct {
	Context *middleware.Context
	Params  DeleteDialParams
	Handler DeleteDialHandler
}

func (o *DeleteDial) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	err := o.Handler.Handle(o.Params) // actually handle the request
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	o.Context.Respond(rw, r, route.Produces, route, nil)

}