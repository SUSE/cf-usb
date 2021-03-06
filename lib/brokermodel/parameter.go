package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
)

/*Parameter A key value parameters

swagger:model Parameter
*/
type Parameter struct {

	/* Name of the parameter
	 */
	Name string `json:"name,omitempty"`

	/* value of the parameter
	 */
	Value interface{} `json:"value,omitempty"`
}

// Validate validates this parameter
func (m *Parameter) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
