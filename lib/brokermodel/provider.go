package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
)

/*Provider Provider

swagger:model Provider
*/
type Provider struct {

	/* Provider Name
	 */
	Name string `json:"name,omitempty"`
}

// Validate validates this provider
func (m *Provider) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
