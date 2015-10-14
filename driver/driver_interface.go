package driver

import "github.com/hpcloud/cf-usb/lib/model"

type Driver interface {
	Init(model.DriverInitRequest, *string) error
	Ping(string, *bool) error
	GetDailsSchema(string, *string) error
	GetConfigSchema(string, *string) error
	ProvisionInstance(model.ProvisionInstanceRequest, *bool) error
	InstanceExists(string, *bool) error
	GenerateCredentials(model.CredentialsRequest, *interface{}) error
	CredentialsExist(model.CredentialsRequest, *bool) error
	RevokeCredentials(model.CredentialsRequest, *interface{}) error
	DeprovisionInstance(string, *interface{}) error
}
