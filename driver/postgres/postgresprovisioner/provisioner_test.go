package postgresprovisioner

import (
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner")

var testPostgresProv = struct {
	postgresProvisioner PostgresProvisionerInterface
	postgresDefaultConn PostgresServiceProperties
}{}

//var userCountQuery = "SELECT COUNT(*) FROM pg_roles WHERE rolname = '%v'"

func init() {
	testPostgresProv.postgresDefaultConn = PostgresServiceProperties{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Dbname:   os.Getenv("POSTGRES_DBNAME"),
		Sslmode:  os.Getenv("POSTGRES_SSLMODE")}

	testPostgresProv.postgresProvisioner = NewPostgresProvisioner(testPostgresProv.postgresDefaultConn, logger)
	testPostgresProv.postgresProvisioner.Init()
}

func TestCreateDatabase(t *testing.T) {
	newDbName := "testcreatedb"

	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'POSTGRES_USER','POSTGRES_PASSWORD','POSTGRES_HOST','POSTGRES_PORT','POSTGRES_DBNAME','POSTGRES_SSLMODE'")
	}

	err := testPostgresProv.postgresProvisioner.CreateDatabase(newDbName)
	if err != nil {
		t.Errorf("Error creating database: ", err)
	}

	exist, err := testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	if err != nil {
		t.Errorf("Error check database exists: ", err)
	}

	if !exist {
		t.Errorf("Database was not created")
	} else {
		t.Log("Database created")
	}
}

func TestCreateUser(t *testing.T) {
	newDbName := "testcreatedb"
	newUser := "testuser"

	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'POSTGRES_USER','POSTGRES_PASSWORD','POSTGRES_HOST','POSTGRES_PORT','POSTGRES_DBNAME','POSTGRES_SSLMODE'")
	}

	exist, err := testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	if err != nil {
		t.Errorf("Error check database exists: ", err)
	}

	if !exist {
		t.Errorf("Database does not exist: ", err)
	}

	err = testPostgresProv.postgresProvisioner.CreateUser(newDbName, newUser, "aPassw0rd")
	if err != nil {
		t.Errorf("Error creating user: ", err)
	}

	exist, err = testPostgresProv.postgresProvisioner.UserExists(newUser)
	if err != nil {
		t.Errorf("Error check user exists: ", err)
	}

	if !exist {
		t.Errorf("User was not created")
	} else {
		t.Log("User created")
	}
}

func TestDeleteUser(t *testing.T) {
	newDbName := "testcreatedb"
	newUser := "testuser"

	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'POSTGRES_USER','POSTGRES_PASSWORD','POSTGRES_HOST','POSTGRES_PORT','POSTGRES_DBNAME','POSTGRES_SSLMODE'")
	}

	exist, err := testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	if err != nil {
		t.Errorf("Error check database exists: ", err)
	}

	if !exist {
		t.Errorf("Database does not exist: ", err)
	}

	exist, err = testPostgresProv.postgresProvisioner.UserExists(newUser)
	if err != nil {
		t.Errorf("Error check user exists: ", err)
	}

	if !exist {
		t.Errorf("User does not exist")
	}

	err = testPostgresProv.postgresProvisioner.DeleteUser(newDbName, newUser)
	if err != nil {
		t.Errorf("Error deleting user: ", err)
	}

	exist, err = testPostgresProv.postgresProvisioner.UserExists(newUser)
	if err != nil {
		t.Errorf("Error check user exists: ", err)
	}

	if !exist {
		t.Log("User was deleted")
	}
}

func TestDeleteDatabase(t *testing.T) {
	newDbName := "testcreatedb"

	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'POSTGRES_USER','POSTGRES_PASSWORD','POSTGRES_HOST','POSTGRES_PORT','POSTGRES_DBNAME','POSTGRES_SSLMODE'")
	}

	exist, err := testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	if err != nil {
		t.Errorf("Error check database exists: ", err)
	}

	if !exist {
		t.Errorf("Database does not exist: ", err)
	}

	err = testPostgresProv.postgresProvisioner.DeleteDatabase(newDbName)
	if err != nil {
		t.Errorf("Error deleting database: ", err)
	}

	exist, err = testPostgresProv.postgresProvisioner.DatabaseExists(newDbName)
	if err != nil {
		t.Errorf("Error check database exists: ", err)
	}

	if !exist {
		t.Log("Database was deleted")
	}
}

func TestParametrizeQuery(t *testing.T) {
	_, err := parametrizeQuery("SELECT COUNT(*) FROM pg_roles WHERE rolname = {{.User}}", map[string]string{"Username": "username"})

	if !strings.Contains(err.Error(), "Invalid parameter passed to query") {
		t.Errorf("Error parametrizing query: ", err)
	}
}

func envVarsOk() bool {
	return testPostgresProv.postgresDefaultConn.User != "" && testPostgresProv.postgresDefaultConn.Password != "" && testPostgresProv.postgresDefaultConn.Host != "" &&
		testPostgresProv.postgresDefaultConn.Port != "" && testPostgresProv.postgresDefaultConn.Dbname != "" && testPostgresProv.postgresDefaultConn.Sslmode != ""
}
