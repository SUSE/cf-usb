package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// GetInfoHandlerFunc turns a function with the right signature into a get info handler
type GetInfoHandlerFunc func(GetInfoParams) (*genmodel.Info, error)

func (fn GetInfoHandlerFunc) Handle(params GetInfoParams) (*genmodel.Info, error) {
	return fn(params)
}

// GetInfoHandler interface for that can handle valid get info params
type GetInfoHandler interface {
	Handle(GetInfoParams) (*genmodel.Info, error)
}

// NewGetInfo creates a new http.Handler for the get info operation
func NewGetInfo(ctx *middleware.Context, handler GetInfoHandler) *GetInfo {
	return &GetInfo{Context: ctx, Handler: handler}
}

/*
Gets information about the USB.

*/
type GetInfo struct {
	Context *middleware.Context
	Params  GetInfoParams
	Handler GetInfoHandler
}

func (o *GetInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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