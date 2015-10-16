package redisprovisioner

type RedisProvisionerInterface interface {
	Init() error
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	GetCredentials(string) (map[string]string, error)
	PingServer() error
}
