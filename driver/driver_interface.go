package driver

import (
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/hpcloud/gocfbroker"
)

type Driver interface {
	Init(config.DriverProperties, *string) error
	Provision(model.DriverProvisionRequest, *string) error
	Deprovision(model.DriverDeprovisionRequest, *string) error
	Update(model.DriverUpdateRequest, *string) error
	Bind(model.DriverBindRequest, *gocfbroker.BindingResponse) error
	Unbind(model.DriverUnbindRequest, *string) error
}
