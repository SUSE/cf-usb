package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// GetInfoHandlerFunc turns a function with the right signature into a get info handler
type GetInfoHandlerFunc func(interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetInfoHandlerFunc) Handle(principal interface{}) middleware.Responder {
	return fn(principal)
}

// GetInfoHandler interface for that can handle valid get info params
type GetInfoHandler interface {
	Handle(interface{}) middleware.Responder
}

// NewGetInfo creates a new http.Handler for the get info operation
func NewGetInfo(ctx *middleware.Context, handler GetInfoHandler) *GetInfo {
	return &GetInfo{Context: ctx, Handler: handler}
}

/*GetInfo swagger:route GET /info getInfo

Gets information about the USB.


*/
type GetInfo struct {
	Context *middleware.Context
	Handler GetInfoHandler
}

func (o *GetInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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

	if err := o.Context.BindValidRequest(r, route, nil); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
