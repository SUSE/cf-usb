package mongoprovisioner

type MongoProvisionerInterface interface {
	IsDatabaseCreated(string) (bool, error)
	IsUserCreated(string, string) (bool, error)
	CreateDatabase(string) error
	DeleteDatabase(string) error
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
}
