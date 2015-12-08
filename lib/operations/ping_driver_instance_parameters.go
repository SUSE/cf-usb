package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/strfmt"
)

// PingDriverInstanceParams contains all the bound params for the ping driver instance operation
// typically these are obtained from a http.Request
//
// swagger:parameters pingDriverInstance
type PingDriverInstanceParams struct {
	/* Driver Instance ID
	Required: true
	In: path
	*/
	DriverInstanceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PingDriverInstanceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	if err := o.bindDriverInstanceID(route.Params.Get("driver_instance_id"), route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PingDriverInstanceParams) bindDriverInstanceID(raw string, formats strfmt.Registry) error {

	o.DriverInstanceID = raw

	return nil
}
