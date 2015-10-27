package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// CreateDialHandlerFunc turns a function with the right signature into a create dial handler
type CreateDialHandlerFunc func(CreateDialParams) error

func (fn CreateDialHandlerFunc) Handle(params CreateDialParams) error {
	return fn(params)
}

// CreateDialHandler interface for that can handle valid create dial params
type CreateDialHandler interface {
	Handle(CreateDialParams) error
}

// NewCreateDial creates a new http.Handler for the create dial operation
func NewCreateDial(ctx *middleware.Context, handler CreateDialHandler) *CreateDial {
	return &CreateDial{Context: ctx, Handler: handler}
}

/*
Create a dial for
*/
type CreateDial struct {
	Context *middleware.Context
	Params  CreateDialParams
	Handler CreateDialHandler
}

func (o *CreateDial) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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