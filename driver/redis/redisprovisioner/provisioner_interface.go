package redisprovisioner

import "github.com/hpcloud/cf-usb/driver/redis/config"

type RedisProvisionerInterface interface {
	Connect(config.RedisDriverConfig) error
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	GetCredentials(string) (map[string]string, error)
	PingServer() error
}
