package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteInstanceParams creates a new DeleteInstanceParams object
// with the default values initialized.
func NewDeleteInstanceParams() DeleteInstanceParams {
	var ()
	return DeleteInstanceParams{}
}

// DeleteInstanceParams contains all the bound params for the delete instance operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteInstance
type DeleteInstanceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*Instance ID
	  Required: true
	  In: path
	*/
	InstanceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *DeleteInstanceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rInstanceID, rhkInstanceID, _ := route.Params.GetOK("instance_id")
	if err := o.bindInstanceID(rInstanceID, rhkInstanceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteInstanceParams) bindInstanceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.InstanceID = raw

	return nil
}