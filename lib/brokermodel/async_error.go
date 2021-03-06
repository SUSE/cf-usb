package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
)

/*AsyncError Async operation not supported error

swagger:model AsyncError
*/
type AsyncError struct {

	/* description
	 */
	Description string `json:"description,omitempty"`

	/* error
	 */
	Error string `json:"error,omitempty"`
}

// Validate validates this async error
func (m *AsyncError) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
