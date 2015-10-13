package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// GetDriversHandlerFunc turns a function with the right signature into a get drivers handler
type GetDriversHandlerFunc func() (*[]string, error)

func (fn GetDriversHandlerFunc) Handle() (*[]string, error) {
	return fn()
}

// GetDriversHandler interface for that can handle valid get drivers params
type GetDriversHandler interface {
	Handle() (*[]string, error)
}

// NewGetDrivers creates a new http.Handler for the get drivers operation
func NewGetDrivers(ctx *middleware.Context, handler GetDriversHandler) *GetDrivers {
	return &GetDrivers{Context: ctx, Handler: handler}
}

/*
Gets information about the available `driver`

*/
type GetDrivers struct {
	Context *middleware.Context
	Handler GetDriversHandler
}

func (o *GetDrivers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	if err := o.Context.BindValidRequest(r, route, nil); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res, err := o.Handler.Handle() // actually handle the request
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	o.Context.Respond(rw, r, route.Produces, route, res)

}
