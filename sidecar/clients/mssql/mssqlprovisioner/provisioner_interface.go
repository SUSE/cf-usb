package mssqlprovisioner

type MssqlProvisionerInterface interface {
	Connect(goSqlDriver string, connectionParams map[string]string) error
	IsDatabaseCreated(databaseId string) (bool, error)
	IsUserCreated(databaseId, userId string) (bool, error)
	CreateDatabase(databaseId string) error
	DeleteDatabase(databaseId string) error
	CreateUser(databaseId, userId, password string) error
	DeleteUser(databaseId, userId string) error
	Close() error
}
