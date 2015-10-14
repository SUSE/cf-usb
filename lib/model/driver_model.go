package model

import "encoding/json"

type ProvisionInstanceRequest struct {
	InstanceID string
	Dails      json.RawMessage
}

type CredentialsRequest struct {
	InstanceID    string
	CredentialsID string
}

type DriverInitRequest struct {
	DriverConfig *json.RawMessage
	Dials        *json.RawMessage
}
