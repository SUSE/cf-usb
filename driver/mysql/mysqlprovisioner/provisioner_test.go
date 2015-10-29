package mysqlprovisioner

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pivotal-golang/lager"
)

var mysqlConConfig = struct {
	User            string
	Pass            string
	Host            string
	TestProvisioner MysqlProvisionerInterface
}{}

//TODO fix tests

func init() {
	var err error
	mysqlConConfig.User = os.Getenv("MYSQL_USER")
	mysqlConConfig.Pass = os.Getenv("MYSQL_PASS")
	mysqlConConfig.Host = os.Getenv("MYSQL_HOST")

	var logger = lager.NewLogger("test-provider")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	mysqlConConfig.TestProvisioner = New(logger)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateDb(t *testing.T) {
	dbName := "test_createdb"
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	log.Println("Creating test database")
	err := mysqlConConfig.TestProvisioner.CreateDatabase(dbName)

	if err != nil {
		log.Fatalln("Error creating database ", err)
	}
}

func TestCreateDbExists(t *testing.T) {
	dbName := "test_createdb"
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}
	log.Println("Testing if database exists")
	created, err := mysqlConConfig.TestProvisioner.IsDatabaseCreated(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if created {
		t.Log("Created true")
	} else {
		t.Log("Created false")
	}
}

func TestCreateUser(t *testing.T) {
	dbName := "test_createdb"

	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	log.Println("Creating test user")
	err := mysqlConConfig.TestProvisioner.CreateUser(dbName, "mytestUser", "mytestPass")
	if err != nil {
		t.Errorf("Error creating user %v", err)
	}
}

func TestCreateUserExists(t *testing.T) {

	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	log.Println("Testing if user exists")
	created, err := mysqlConConfig.TestProvisioner.IsUserCreated("mytestUser")
	if err != nil {
		t.Errorf("Error verifying user %v", err)
	}
	if created {
		t.Log("test user is created")
	} else {
		t.Log("test user was not created")
	}
}

func TestDeleteUser(t *testing.T) {
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	log.Println("Removing test user")
	err := mysqlConConfig.TestProvisioner.DeleteUser("mytestUser")
	if err != nil {
		t.Errorf("Error deleting user %v", err)
	}
}

func TestDeleteTheDatabase(t *testing.T) {
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	dbName := "test_createdb"
	log.Println("Removing test database")

	err := mysqlConConfig.TestProvisioner.DeleteDatabase(dbName)
	if err != nil {
		t.Errorf("Error deleting database %v", err)
	}
}
