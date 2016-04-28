package genmodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

/*Instance instance

swagger:model instance
*/
type Instance struct {

	/* configuration
	 */
	Configuration interface{} `json:"configuration,omitempty"`

	/* dials
	 */
	Dials []string `json:"dials,omitempty"`

	/* id
	 */
	ID string `json:"id,omitempty"`

	/* name

	Required: true
	*/
	Name *string `json:"name"`

	/* service
	 */
	Service string `json:"service,omitempty"`

	/* target URL
	 */
	TargetURL string `json:"targetURL,omitempty"`
}

// Validate validates this instance
func (m *Instance) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDials(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Instance) validateDials(formats strfmt.Registry) error {

	if swag.IsZero(m.Dials) { // not required
		return nil
	}

	return nil
}

func (m *Instance) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}
