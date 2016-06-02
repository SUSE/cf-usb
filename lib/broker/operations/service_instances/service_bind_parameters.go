package service_instances

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/hpcloud/cf-usb/lib/brokermodel"
)

// NewServiceBindParams creates a new ServiceBindParams object
// with the default values initialized.
func NewServiceBindParams() ServiceBindParams {
	var ()
	return ServiceBindParams{}
}

// ServiceBindParams contains all the bound params for the service bind operation
// typically these are obtained from a http.Request
//
// swagger:parameters serviceBind
type ServiceBindParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*
	  Required: true
	  In: body
	*/
	Binding *brokermodel.Binding
	/*The binding_id of a service binding is provided by the Cloud Controller.
	  Required: true
	  In: path
	*/
	BindingID string
	/*The instance_id of a service instance is provided by the Cloud Controller. This ID will be used for future requests (bind and deprovision), so the broker must use it to correlate the resource it creates.
	  Required: true
	  In: path
	*/
	InstanceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *ServiceBindParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body brokermodel.Binding
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("binding", "body"))
			} else {
				res = append(res, errors.NewParseError("binding", "body", "", err))
			}

		} else {
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Binding = &body
			}
		}

	} else {
		res = append(res, errors.Required("binding", "body"))
	}

	rBindingID, rhkBindingID, _ := route.Params.GetOK("binding_id")
	if err := o.bindBindingID(rBindingID, rhkBindingID, route.Formats); err != nil {
		res = append(res, err)
	}

	rInstanceID, rhkInstanceID, _ := route.Params.GetOK("instance_id")
	if err := o.bindInstanceID(rInstanceID, rhkInstanceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ServiceBindParams) bindBindingID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.BindingID = raw

	return nil
}

func (o *ServiceBindParams) bindInstanceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.InstanceID = raw

	return nil
}
