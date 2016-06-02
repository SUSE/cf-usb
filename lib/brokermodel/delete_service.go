package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
)

/*DeleteService Parameters needed to unbind a service instance

swagger:model DeleteService
*/
type DeleteService struct {

	/* A value of true indicates that both the Cloud Controller and the requesting client support asynchronous provisioning. If this parameter is not included in the request, and the broker can only provision an instance of the requested plan asynchronously, the broker should reject the request with a 422 as described below
	 */
	AcceptsIncomplete bool `json:"accepts_incomplete,omitempty"`

	/* ID of the plan from the catalog. While not strictly necessary, some brokers might make use of this ID.
	 */
	PlanID string `json:"plan_id,omitempty"`

	/* ID of the service from the catalog. While not strictly necessary, some brokers might make use of this ID.
	 */
	ServiceID string `json:"service_id,omitempty"`
}

// Validate validates this delete service
func (m *DeleteService) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
