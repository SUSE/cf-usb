package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
)

// CreateServicePlanHandlerFunc turns a function with the right signature into a create service plan handler
type CreateServicePlanHandlerFunc func(CreateServicePlanParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateServicePlanHandlerFunc) Handle(params CreateServicePlanParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CreateServicePlanHandler interface for that can handle valid create service plan params
type CreateServicePlanHandler interface {
	Handle(CreateServicePlanParams, interface{}) middleware.Responder
}

// NewCreateServicePlan creates a new http.Handler for the create service plan operation
func NewCreateServicePlan(ctx *middleware.Context, handler CreateServicePlanHandler) *CreateServicePlan {
	return &CreateServicePlan{Context: ctx, Handler: handler}
}

/*CreateServicePlan swagger:route POST /plans createServicePlan

Create a plan

*/
type CreateServicePlan struct {
	Context *middleware.Context
	Params  CreateServicePlanParams
	Handler CreateServicePlanHandler
}

func (o *CreateServicePlan) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
