package rabbitmqprovisioner

import "github.com/hpcloud/cf-usb/driver/rabbitmq/config"

type RabbitmqProvisionerInterface interface {
	Connect(conf config.RabbitmqDriverConfig) error
	CreateDatabase(string) error
	DeleteDatabase(string) error
	DatabaseExists(string) (bool, error)
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
	UserExists(string) (bool, error)
}
