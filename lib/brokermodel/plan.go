package brokermodel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

/*Plan A plan for the service

swagger:model Plan
*/
type Plan struct {

	/* A short description of the service that will appear in the catalog.
	 */
	Description string `json:"description,omitempty"`

	/* This field allows the plan to be limited by the non_basic_services_allowed field in a Cloud Foundry Quota, [see Quota Plans](http://docs.cloudfoundry.org/running/managing-cf/quota-plans.html).
	 */
	Free bool `json:"free,omitempty"`

	/* An identifier used to correlate this plan in future requests to the catalog. This must be unique within Cloud Foundry, using a GUID is recommended.
	 */
	ID string `json:"id,omitempty"`

	/* metadata
	 */
	Metadata *PlanMetadata `json:"metadata,omitempty"`

	/* The CLI-friendly name of the plan that will appear in the catalog. All lowercase, no spaces.
	 */
	Name string `json:"name,omitempty"`
}

// Validate validates this plan
func (m *Plan) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetadata(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Plan) validateMetadata(formats strfmt.Registry) error {

	if swag.IsZero(m.Metadata) { // not required
		return nil
	}

	if m.Metadata != nil {

		if err := m.Metadata.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}

/*PlanMetadata A list of metadata for a service plan. For more information, [see Service Metadata](https://docs.cloudfoundry.org/services/catalog-metadata.html).

swagger:model PlanMetadata
*/
type PlanMetadata struct {

	/* A description of the service plan to be displayed in a catalog.
	 */
	Description string `json:"description,omitempty"`

	/* Additional non mandatory fields for Plan metadata (e.g. metadata.displayName, metadata.bullets)
	 */
	Metadata interface{} `json:"metadata,omitempty"`

	/* A short name for the service plan to be displayed in a catalog.
	 */
	Name string `json:"name,omitempty"`
}

// Validate validates this plan metadata
func (m *PlanMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
