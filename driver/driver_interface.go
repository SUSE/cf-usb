package driver

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
)

type Driver interface {
	Init(config.DriverProperties, *string) error
	Provision(model.DriverProvisionRequest, *string) error
	Deprovision(model.DriverDeprovisionRequest, *string) error
	Bind(model.DriverBindRequest, *json.RawMessage) error
	Unbind(model.DriverUnbindRequest, *string) error
}
