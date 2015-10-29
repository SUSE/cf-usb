package postgresprovisioner

import "github.com/hpcloud/cf-usb/driver/postgres/config"

type PostgresProvisionerInterface interface {
	Connect(conf config.PostgresDriverConfig) error
	CreateDatabase(string) error
	DeleteDatabase(string) error
	DatabaseExists(string) (bool, error)
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
	UserExists(string) (bool, error)
}
