package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/strfmt"

	"github.com/hpcloud/cf-usb/lib/genmodel"
)

// UpdateServiceParams contains all the bound params for the update service operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateService
type UpdateServiceParams struct {
	/* Update service
	Required: true
	In: body
	*/
	Service *genmodel.Service
	/* ID of the service
	Required: true
	In: path
	*/
	ServiceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *UpdateServiceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	var body genmodel.Service
	if err := route.Consumer.Consume(r.Body, &body); err != nil {
		res = append(res, errors.NewParseError("service", "body", "", err))
	} else {
		if err := body.Validate(route.Formats); err != nil {
			res = append(res, err)
		}

		if len(res) == 0 {
			o.Service = &body
		}
	}

	if err := o.bindServiceID(route.Params.Get("service_id"), route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateServiceParams) bindServiceID(raw string, formats strfmt.Registry) error {

	o.ServiceID = raw

	return nil
}
