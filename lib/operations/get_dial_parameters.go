package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/strfmt"
)

// NewGetDialParams creates a new GetDialParams object
// with the default values initialized.
func NewGetDialParams() GetDialParams {
	var ()
	return GetDialParams{}
}

// GetDialParams contains all the bound params for the get dial operation
// typically these are obtained from a http.Request
//
// swagger:parameters getDial
type GetDialParams struct {
	/*ID of the dial
	  Required: true
	  In: path
	*/
	DialID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetDialParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	rDialID, rhkDialID, _ := route.Params.GetOK("dial_id")
	if err := o.bindDialID(rDialID, rhkDialID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDialParams) bindDialID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.DialID = raw

	return nil
}
