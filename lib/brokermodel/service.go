package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

/*Service Service object

swagger:model Service
*/
type Service struct {

	/* A value of true indicates that both the Cloud Controller and the requesting client support asynchronous provisioning. If this parameter is not included in the request, and the broker can only provision an instance of the requested plan asynchronously, the broker should reject the request with a 422 as described below
	 */
	AcceptsIncomplete bool `json:"accepts_incomplete,omitempty"`

	/* The Cloud Controller GUID of the organization under which the service is to be provisioned. Although most brokers will not use this field, it could be helpful in determining data placement or applying custom business rules.
	 */
	OrganizationGUID string `json:"organization_guid,omitempty"`

	/* parameteres
	 */
	Parameteres *Parameter `json:"parameteres,omitempty"`

	/* The ID of the plan within the above service (from the catalog endpoint) that the user would like provisioned. Because plans have identifiers unique to a broker, this is enough information to determine what to provision.
	 */
	PlanID string `json:"plan_id,omitempty"`

	/* The ID of the service within the catalog above. While not strictly necessary, some brokers might make use of this ID.
	 */
	ServiceID string `json:"service_id,omitempty"`

	/* Similar to organization_guid, but for the space.
	 */
	SpaceGUID string `json:"space_guid,omitempty"`
}

// Validate validates this service
func (m *Service) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateParameteres(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Service) validateParameteres(formats strfmt.Registry) error {

	if swag.IsZero(m.Parameteres) { // not required
		return nil
	}

	if m.Parameteres != nil {

		if err := m.Parameteres.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}