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
