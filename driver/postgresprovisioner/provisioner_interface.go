package postgresprovisioner

type PostgresProvisionerInterface interface {
	CreateDatabase(string) error
	DeleteDatabase(string) error
	CreateUser(string, string, string) error
	DeleteUser(string, string) error
}
