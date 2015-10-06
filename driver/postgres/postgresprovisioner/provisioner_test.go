package postgresprovisioner

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pivotal-golang/lager/lagertest"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner")

var postgresDefaultConn = PostgresServiceProperties{User: "postgres", Password: "password1234!", Host: "localhost", Port: "5432", Dbname: "postgres", Sslmode: "disable"}

var userCountQuery = "SELECT COUNT(*) FROM pg_roles WHERE rolname = '%v'"

func TestCreateDatabase(t *testing.T) {
	newDbName := "testcreatedb"
	fmt.Println("conn string: ", postgresDefaultConn)

	testp := NewPostgresProvisioner(postgresDefaultConn, logger)
	testp.Init()

	err := testp.CreateDatabase(newDbName)
	if err != nil {
		t.Errorf("Error creating database: ", err)
	}
}

func TestCreateUser(t *testing.T) {
	newDbName := "testcreatedb"
	newUser := "testuser"

	testp := NewPostgresProvisioner(postgresDefaultConn, logger)
	testp.Init()

	err := testp.CreateUser(newDbName, newUser, "aPassw0rd")
	if err != nil {
		t.Errorf("Error creating user: ", err)
	}

	pgClient, err := sql.Open("postgres", buildConnectionString(postgresDefaultConn))
	if err != nil {
		t.Errorf("Error opening postgres client: ", err)
	}
	defer pgClient.Close()

	userCount := 0
	err = pgClient.QueryRow(fmt.Sprintf(userCountQuery, newUser)).Scan(&userCount)
	if err != nil {
		t.Errorf("Error executing query: ", err)
	}

	if userCount == 0 {
		t.Errorf("User was not created: ", err)
	}
}

func TestDeleteUser(t *testing.T) {
	newDbName := "testcreatedb"
	newUser := "testuser"

	testp := NewPostgresProvisioner(postgresDefaultConn, logger)
	testp.Init()

	err := testp.DeleteUser(newDbName, newUser)
	if err != nil {
		t.Errorf("Error deleting user: ", err)
	}

	pgClient, err := sql.Open("postgres", buildConnectionString(postgresDefaultConn))
	if err != nil {
		t.Errorf("Error opening postgres client: ", err)
	}
	defer pgClient.Close()

	userCount := 0
	err = pgClient.QueryRow(fmt.Sprintf(userCountQuery, newUser)).Scan(&userCount)
	if err != nil {
		t.Errorf("Error executing query: ", err)
	}

	if userCount > 0 {
		t.Errorf("User was not created: ", err)
	}
}

func TestDeleteDatabase(t *testing.T) {
	newDbName := "testcreatedb"

	testp := NewPostgresProvisioner(postgresDefaultConn, logger)
	testp.Init()

	err := testp.DeleteDatabase(newDbName)
	if err != nil {
		t.Errorf("Error deleting database: ", err)
	}
}

func TestParametrizeQuery(t *testing.T) {
	_, err := parametrizeQuery("SELECT COUNT(*) FROM pg_roles WHERE rolname = {{.User}}", map[string]string{"Username": "username"})

	if !strings.Contains(err.Error(), "Invalid parameter passed to query") {
		t.Errorf("Error parametrizing query: ", err)
	}
}
