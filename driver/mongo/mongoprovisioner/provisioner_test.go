package mongoprovisioner

import (
	"log"
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/driver/mongo/config"
	"github.com/pivotal-golang/lager"
)

var mongoConConfig = struct {
	User            string
	Pass            string
	Host            string
	Port            string
	TestProvisioner MongoProvisionerInterface
}{}

func init() {
	var err error
	mongoConConfig.User = os.Getenv("MONGO_USER")
	mongoConConfig.Pass = os.Getenv("MONGO_PASS")
	mongoConConfig.Host = os.Getenv("MONGO_HOST")
	mongoConConfig.Port = os.Getenv("MONGO_PORT")

	var logger = lager.NewLogger("test-mongo-provider")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	mongoConConfig.TestProvisioner = New(logger)

	mongo := &config.MongoDriverConfig{
		Host: mongoConConfig.Host,
		Port: mongoConConfig.Port,
		Pass: mongoConConfig.Pass,
		User: mongoConConfig.User,
	}

	err = mongoConConfig.TestProvisioner.Connect(*mongo)
	if err != nil {
		log.Println(err)
	}
}

func TestCreateDb(t *testing.T) {
	dbName := "test_createdb"
	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Creating test database")
	err := mongoConConfig.TestProvisioner.CreateDatabase(dbName)

	if err != nil {
		log.Fatalln("Error creating database ", err)
	}
}

func TestCreateDbExists(t *testing.T) {
	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Testing if database exists")
	created, err := mongoConConfig.TestProvisioner.IsDatabaseCreated(dbName)
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

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Creating test user")
	err := mongoConConfig.TestProvisioner.CreateUser(dbName, "mytestUser", "mytestPass")
	if err != nil {
		t.Errorf("Error creating user %v", err)
	}
}

func TestCreateUserExists(t *testing.T) {
	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Testing if user exists")
	created, err := mongoConConfig.TestProvisioner.IsUserCreated(dbName, "mytestUser")
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
	dbName := "test_createdb"

	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	log.Println("Removing test user")
	err := mongoConConfig.TestProvisioner.DeleteUser(dbName, "mytestUser")
	if err != nil {
		t.Errorf("Error deleting user %v", err)
	}
}

func TestDeleteTheDatabase(t *testing.T) {
	if mongoConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MONGO_USER','MONGO_PASS','MONGO_HOST', 'MONGO_PORT'")
	}

	dbName := "test_createdb"
	log.Println("Removing test database")

	err := mongoConConfig.TestProvisioner.DeleteDatabase(dbName)
	if err != nil {
		t.Errorf("Error deleting database %v", err)
	}
}
