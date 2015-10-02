package mysqlprovisioner

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var mysqlConConfig = struct {
	User string
	Pass string
	Host string
}{
	User: os.Getenv("MYSQL_USER"),
	Pass: os.Getenv("MYSQL_PASSWORD"),
	Host: os.Getenv("MYSQL_HOST"),
}

func TestCreateDb(t *testing.T) {
	dbName := "test_createdb"
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)
	log.Println("Creating test database")
	err := tp.CreateDatabase(dbName)

	if err != nil {
		log.Fatalln("Error creating database ", err)
	}
}

func TestCreateUser(t *testing.T) {
	dbName := "test_createdb"

	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)
	log.Println("Creating test user")
	err := tp.CreateUser(dbName, "mytestUser", "mytestPass")
	if err != nil {
		t.Errorf("Error creating user %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	if mysqlConConfig.User == "" || mysqlConConfig.Pass == "" || mysqlConConfig.Host == "" {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST'(ip:port)")
	}

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)
	log.Println("Removing test user")
	err := tp.DeleteUser("mytestUser")
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
	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)

	err := tp.DeleteDatabase(dbName)
	if err != nil {
		t.Errorf("Error deleting database %v", err)
	}
}
