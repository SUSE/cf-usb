package model

import "github.com/pivotal-cf/brokerapi"

type DriverProvisionRequest struct {
	InstanceID     string
	ServiceDetails brokerapi.ProvisionDetails
}

type DriverDeprovisionRequest struct {
	InstanceID string
}

type DriverBindRequest struct {
	InstanceID  string
	BindingID   string
	BindDetails brokerapi.BindDetails
}

type DriverUnbindRequest struct {
	InstanceID string
	BindingID  string
}
