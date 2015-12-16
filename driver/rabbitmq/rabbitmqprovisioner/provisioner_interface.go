package rabbitmqprovisioner

import "github.com/hpcloud/cf-usb/driver/rabbitmq/config"

type RabbitmqProvisionerInterface interface {
	Connect(config.RabbitmqDriverConfig) error
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	CreateUser(string, string) (map[string]string, error)
	DeleteUser(string, string) error
	UserExists(string, string) (bool, error)
	PingServer() error
}
