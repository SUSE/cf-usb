package mysqlprovisioner

import (
	"database/sql"

	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/config"
)

type MysqlProvisionerInterface interface {
	Connect(config.MysqlDriverConfig) error
	IsDatabaseCreated(string) (bool, error)
	IsUserCreated(string) (bool, error)
	CreateDatabase(string) error
	DeleteDatabase(string) error
	Query(string, ...interface{}) (*sql.Rows, error)
	CreateUser(string, string, string) error
	DeleteUser(string) error
}
