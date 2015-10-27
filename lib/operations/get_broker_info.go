package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// GetBrokerInfoHandlerFunc turns a function with the right signature into a get broker info handler
type GetBrokerInfoHandlerFunc func(GetBrokerInfoParams) (*genmodel.Broker, error)

func (fn GetBrokerInfoHandlerFunc) Handle(params GetBrokerInfoParams) (*genmodel.Broker, error) {
	return fn(params)
}

// GetBrokerInfoHandler interface for that can handle valid get broker info params
type GetBrokerInfoHandler interface {
	Handle(GetBrokerInfoParams) (*genmodel.Broker, error)
}

// NewGetBrokerInfo creates a new http.Handler for the get broker info operation
func NewGetBrokerInfo(ctx *middleware.Context, handler GetBrokerInfoHandler) *GetBrokerInfo {
	return &GetBrokerInfo{Context: ctx, Handler: handler}
}

/*
Gets the broker api connection info

*/
type GetBrokerInfo struct {
	Context *middleware.Context
	Params  GetBrokerInfoParams
	Handler GetBrokerInfoHandler
}

func (o *GetBrokerInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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