package model

import "github.com/hpcloud/gocfbroker"

type DriverProvisionRequest struct {
	InstanceID             string
	BrokerProvisionRequest gocfbroker.ProvisionRequest
}

type DriverDeprovisionRequest struct {
	InstanceID string
	PlanID     string
}

type DriverUpdateRequest struct {
	InstanceID          string
	BrokerUpdateRequest gocfbroker.UpdateProvisionRequest
}

type DriverBindRequest struct {
	InstanceID        string
	BindingID         string
	BrokerBindRequest gocfbroker.BindingRequest
}

type DriverUnbindRequest struct {
	InstanceID string
	BindingID  string
	PlanID     string
}
