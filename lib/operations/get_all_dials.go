package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// GetAllDialsHandlerFunc turns a function with the right signature into a get all dials handler
type GetAllDialsHandlerFunc func(GetAllDialsParams) (*[]genmodel.Dial, error)

func (fn GetAllDialsHandlerFunc) Handle(params GetAllDialsParams) (*[]genmodel.Dial, error) {
	return fn(params)
}

// GetAllDialsHandler interface for that can handle valid get all dials params
type GetAllDialsHandler interface {
	Handle(GetAllDialsParams) (*[]genmodel.Dial, error)
}

// NewGetAllDials creates a new http.Handler for the get all dials operation
func NewGetAllDials(ctx *middleware.Context, handler GetAllDialsHandler) GetAllDials {
	return GetAllDials{Context: ctx, Handler: handler}
}

/*
Gets `dials`
*/
type GetAllDials struct {
	Context *middleware.Context
	Params  GetAllDialsParams
	Handler GetAllDialsHandler
}

func (o GetAllDials) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
