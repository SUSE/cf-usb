package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// UpdateDialHandlerFunc turns a function with the right signature into a update dial handler
type UpdateDialHandlerFunc func(UpdateDialParams) error

func (fn UpdateDialHandlerFunc) Handle(params UpdateDialParams) error {
	return fn(params)
}

// UpdateDialHandler interface for that can handle valid update dial params
type UpdateDialHandler interface {
	Handle(UpdateDialParams) error
}

// NewUpdateDial creates a new http.Handler for the update dial operation
func NewUpdateDial(ctx *middleware.Context, handler UpdateDialHandler) *UpdateDial {
	return &UpdateDial{Context: ctx, Handler: handler}
}

/*
Updates the dial with the id **dial_id**
*/
type UpdateDial struct {
	Context *middleware.Context
	Params  UpdateDialParams
	Handler UpdateDialHandler
}

func (o *UpdateDial) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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