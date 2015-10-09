package mysqlprovisioner

import (
	"database/sql"
)

type MysqlProvisionerInterface interface {
	IsDatabaseCreated(string) (bool, error)
	IsUserCreated(string) (bool, error)
	CreateDatabase(string) error
	DeleteDatabase(string) error
	Query(string) (*sql.Rows, error)
	CreateUser(string, string, string) error
	DeleteUser(string) error
}
