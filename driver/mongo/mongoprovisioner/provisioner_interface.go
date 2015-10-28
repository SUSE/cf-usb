package mongoprovisioner

import "github.com/hpcloud/cf-usb/driver/mongo/config"

type MongoProvisionerInterface interface {
	Connect(config.MongoDriverConfig) error
	IsDatabaseCreated(string) (bool, error)
	IsUserCreated(string, string) (bool, error)
	CreateDatabase(string) error
	DeleteDatabase(string) error
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
}
