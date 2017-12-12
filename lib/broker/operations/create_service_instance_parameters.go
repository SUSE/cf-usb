package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/SUSE/cf-usb/lib/brokermodel"
)

// NewCreateServiceInstanceParams creates a new CreateServiceInstanceParams object
// with the default values initialized.
func NewCreateServiceInstanceParams() CreateServiceInstanceParams {
	var ()
	return CreateServiceInstanceParams{}
}

// CreateServiceInstanceParams contains all the bound params for the create service instance operation
// typically these are obtained from a http.Request
//
// swagger:parameters createServiceInstance
type CreateServiceInstanceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*The instance_id of a service instance is provided by the Cloud Controller. This ID will be used for future requests (bind and deprovision), so the broker must use it to correlate the resource it creates.
	  Required: true
	  In: path
	*/
	InstanceID string
	/*Service information.
	  Required: true
	  In: body
	*/
	Service *brokermodel.Service
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *CreateServiceInstanceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rInstanceID, rhkInstanceID, _ := route.Params.GetOK("instance_id")
	if err := o.bindInstanceID(rInstanceID, rhkInstanceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body brokermodel.Service
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("service", "body"))
			} else {
				res = append(res, errors.NewParseError("service", "body", "", err))
			}

		} else {
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Service = &body
			}
		}

	} else {
		res = append(res, errors.Required("service", "body"))
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateServiceInstanceParams) bindInstanceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.InstanceID = raw

	return nil
}
