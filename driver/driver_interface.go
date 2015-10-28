package driver

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver/status"
)

type ProvisionInstanceRequest struct {
	InstanceID string
	Config     *json.RawMessage
	Dials      *json.RawMessage
}

type GetInstanceRequest struct {
	InstanceID string
	Config     *json.RawMessage
}

type DeprovisionInstanceRequest struct {
	InstanceID string
	Config     *json.RawMessage
}

type GenerateCredentialsRequest struct {
	InstanceID    string
	CredentialsID string
	Config        *json.RawMessage
}

type GetCredentialsRequest struct {
	InstanceID    string
	CredentialsID string
	Config        *json.RawMessage
}

type RevokeCredentialsRequest struct {
	InstanceID    string
	CredentialsID string
	Config        *json.RawMessage
}

type Instance struct {
	InstanceID  string
	Status      status.Status
	Description string
}

type Credentials struct {
	CredentialsID string
	Status        status.Status
	Description   string
}

type Driver interface {
	Ping(*json.RawMessage, *bool) error
	GetDailsSchema(string, *string) error
	GetConfigSchema(string, *string) error
	ProvisionInstance(ProvisionInstanceRequest, *Instance) error
	GetInstance(GetInstanceRequest, *Instance) error
	GenerateCredentials(GenerateCredentialsRequest, *interface{}) error
	GetCredentials(GetCredentialsRequest, *Credentials) error
	RevokeCredentials(RevokeCredentialsRequest, *Credentials) error
	DeprovisionInstance(DeprovisionInstanceRequest, *Instance) error
}
