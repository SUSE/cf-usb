package gocfbroker

import (
	"bytes"
	"encoding/json"

	"github.com/kat-co/vala"
)

const (
	// brokerJSONSchemaVersion is the version of the current broker's JSON
	// schema.
	brokerJSONSchemaVersion = 1
)

// Instance is stored inside the database for each service instance.
type Instance struct {
	ProvisionReq ProvisionRequest  `json:"provision_request"`
	ProvisionRes ProvisionResponse `json:"provision_response"`

	Bindings []Binding `json:"bindings"`

	// Version is the json schema version number
	Version int64 `json:"version"`
}

// Binding is part of the Instance storage.
type Binding struct {
	BindingID  string         `json:"instance_id"`
	BindingReq BindingRequest `json:"binding_request"`
	BindingRes string         `json:"binding_response"`
}

// ProvisionRequest to try to provision a service.
type ProvisionRequest struct {
	ServiceID        string           `json:"service_id"`
	PlanID           string           `json:"plan_id"`
	OrganizationGUID string           `json:"organization_guid"`
	SpaceGUID        string           `json:"space_guid"`
	Parameters       *json.RawMessage `json:"parameters,omitempty"`
}

// validate the request.
func (p ProvisionRequest) validate(c Catalog) *vala.Validation {
	return vala.BeginValidation().Validate(
		validateServiceID(c, p.ServiceID),
		validatePlanID(c, p.ServiceID, p.PlanID),
		validateStrNotEmpty(p.OrganizationGUID, "organization_guid"),
		validateStrNotEmpty(p.SpaceGUID, "space_guid"),
	)
}

// Equal returns true of the two provision requests are equivalent.
func (p ProvisionRequest) Equal(other ProvisionRequest) bool {
	if p.ServiceID != other.ServiceID {
		return false
	}
	if p.PlanID != other.PlanID {
		return false
	}
	if p.OrganizationGUID != other.OrganizationGUID {
		return false
	}
	if p.SpaceGUID != other.SpaceGUID {
		return false
	}

	return jsonRawEqual(p.Parameters, other.Parameters)
}

// ProvisionResponse is returned in response to ProvisionRequests.
type ProvisionResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

// UpdateProvisionRequest updates a service with a new plan.
type UpdateProvisionRequest struct {
	ServiceID      string           `json:"service_id"`
	PlanID         string           `json:"plan_id"`
	PreviousValues *PreviousValues  `json:"previous_values,omitempty"`
	Parameters     *json.RawMessage `json:"parameters,omitempty"`
}

// PreviousValues for the UpdateProvisionRequest
type PreviousValues struct {
	ServiceID      string `json:"service_id,omitempty"`
	PlanID         string `json:"plan_id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	SpaceID        string `json:"space_id,omitempty"`
}

// validate the updateProvisionRequest
func (u UpdateProvisionRequest) validate(c Catalog) *vala.Validation {
	return vala.BeginValidation().Validate(
		validateServiceID(c, u.ServiceID),
		validatePlanID(c, u.ServiceID, u.PlanID),
		validateServiceUpdatable(c, u.ServiceID),
	)
}

// Equal returns true of the two update provision requests are equivalent.
func (u UpdateProvisionRequest) Equal(other UpdateProvisionRequest) bool {
	if u.PlanID != other.PlanID {
		return false
	}

	if u.Parameters == nil && other.Parameters == nil {
		return true
	}

	return jsonRawEqual(u.Parameters, other.Parameters)
}

// BindingRequest to bind a service to an application.
type BindingRequest struct {
	ServiceID  string           `json:"service_id"`
	PlanID     string           `json:"plan_id"`
	AppGUID    string           `json:"app_guid,omitempty"`
	Parameters *json.RawMessage `json:"parameters,omitempty"`
}

// validate the request.
func (b BindingRequest) validate(c Catalog) *vala.Validation {
	return vala.BeginValidation().Validate(
		validateServiceID(c, b.ServiceID),
		validatePlanID(c, b.ServiceID, b.PlanID),
	)
}

// Equal returns true of the two binding requests are equivalent.
func (b BindingRequest) Equal(other BindingRequest) bool {
	if b.ServiceID != other.ServiceID {
		return false
	}
	if b.PlanID != other.PlanID {
		return false
	}
	if b.AppGUID != other.AppGUID {
		return false
	}

	return jsonRawEqual(b.Parameters, other.Parameters)
}

// BindingResponse for BindingRequests.
type BindingResponse struct {
	Credentials    *json.RawMessage `json:"credentials"`
	SyslogDrainURL string           `json:"syslog_drain_url,omitempty"`
}

// Equal returs true if the two binding responses are equivalent.
func (b BindingResponse) Equal(other BindingResponse) bool {
	if b.SyslogDrainURL != other.SyslogDrainURL {
		return false
	}

	return jsonRawEqual(b.Credentials, other.Credentials)
}

// Options is the base configuration required to start a v2 broker service
// and advertise its offering.
type Options struct {
	APIVersion     string `json:"api_version"`
	AuthUser       string `json:"auth_user"`
	AuthPassword   string `json:"auth_password"`
	RequireAppGUID bool   `json:"require_app_guid_in_bind_requests"`
	Listen         string `json:"listen"`
	EncryptionKey  string `json:"db_encryption_key"`
	Catalog
}

// Catalog provides the list of services the service offers.
type Catalog struct {
	Services []Service `json:"services"`
}

// Service describes the service in its entirety.
type Service struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Bindable      bool             `json:"bindable"`
	Tags          []string         `json:"tags,omitempty"`
	Requires      []string         `json:"requires,omitempty"`
	PlanUpdatable bool             `json:"plan_updateable"`
	Plans         []Plan           `json:"plans"`
	Metadata      *json.RawMessage `json:"metadata,omitempty"`
}

// Plan describes the type of plan and the quotas associated.
type Plan struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Free        bool             `json:"free,omitempty"`
	Metadata    *json.RawMessage `json:"metadata,omitempty"`
}

// DashboardClient contains the data necessary to activate the Dashboard SSO
// feature for this service.
type DashboardClient struct {
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// deprovisionRequest is used internally for job queueing
type deprovisionRequest struct {
	serviceID, planID string
}

// MakeJSONRawMessage is a helper to convert a string into the required
// pointer type for json.RawMessage serialization to work properly.
func MakeJSONRawMessage(str string) *json.RawMessage {
	j := json.RawMessage(str)
	return &j
}

// jsonRawEqual checks to see if two pointers to []byte are equal.
func jsonRawEqual(a, b *json.RawMessage) bool {
	if a == nil && b == nil {
		return true
	} else if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return bytes.Compare(*a, *b) == 0
}
