package mysqlprovisioner

import (
	"database/sql"
)

type MysqlProvisionerInterface interface {
	CreateDatabase(string) error
	DeleteDatabase(string) error
	Query(string) (*sql.Rows, error)
	CreateUser(string, string, string) error
	DeleteUser(string) error
}
