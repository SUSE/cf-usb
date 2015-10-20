package postgresprovisioner

type PostgresProvisionerInterface interface {
	Init() error
	Ping() error
	CreateDatabase(string) error
	DeleteDatabase(string) error
	DatabaseExists(string) (bool, error)
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
	UserExists(string) (bool, error)
}
