package driver

import "github.com/hpcloud/cf-usb/lib/config"

type Driver interface {
	Init(config.DriverProperties, *string) error
	Provision(interface{}, *interface{}) error
	Deprovision(string, *string) error
	Bind(string, *string) error
	Unbind(string, *string) error
	Update(string, *string) error
}
