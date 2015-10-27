package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// GetDriverSchemaHandlerFunc turns a function with the right signature into a get driver schema handler
type GetDriverSchemaHandlerFunc func(GetDriverSchemaParams) (*genmodel.DriverSchema, error)

func (fn GetDriverSchemaHandlerFunc) Handle(params GetDriverSchemaParams) (*genmodel.DriverSchema, error) {
	return fn(params)
}

// GetDriverSchemaHandler interface for that can handle valid get driver schema params
type GetDriverSchemaHandler interface {
	Handle(GetDriverSchemaParams) (*genmodel.DriverSchema, error)
}

// NewGetDriverSchema creates a new http.Handler for the get driver schema operation
func NewGetDriverSchema(ctx *middleware.Context, handler GetDriverSchemaHandler) *GetDriverSchema {
	return &GetDriverSchema{Context: ctx, Handler: handler}
}

/*
Get driver config schema
*/
type GetDriverSchema struct {
	Context *middleware.Context
	Params  GetDriverSchemaParams
	Handler GetDriverSchemaHandler
}

func (o *GetDriverSchema) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res, err := o.Handler.Handle(o.Params) // actually handle the request
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	o.Context.Respond(rw, r, route.Produces, route, res)

}