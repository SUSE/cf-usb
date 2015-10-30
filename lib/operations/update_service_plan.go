package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// UpdateServicePlanHandlerFunc turns a function with the right signature into a update service plan handler
type UpdateServicePlanHandlerFunc func(UpdateServicePlanParams) (*genmodel.Plan, error)

func (fn UpdateServicePlanHandlerFunc) Handle(params UpdateServicePlanParams) (*genmodel.Plan, error) {
	return fn(params)
}

// UpdateServicePlanHandler interface for that can handle valid update service plan params
type UpdateServicePlanHandler interface {
	Handle(UpdateServicePlanParams) (*genmodel.Plan, error)
}

// NewUpdateServicePlan creates a new http.Handler for the update service plan operation
func NewUpdateServicePlan(ctx *middleware.Context, handler UpdateServicePlanHandler) *UpdateServicePlan {
	return &UpdateServicePlan{Context: ctx, Handler: handler}
}

/*
Updates the plan with the id **planID** for the service id **serviceID**
*/
type UpdateServicePlan struct {
	Context *middleware.Context
	Params  UpdateServicePlanParams
	Handler UpdateServicePlanHandler
}

func (o *UpdateServicePlan) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
