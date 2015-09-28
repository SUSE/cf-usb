package mysqlprovisioner

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

var mysqlConConfig = struct {
	User string
	Pass string
	Host string
}{
	User: "root",
	Pass: "password1234",
	Host: "127.0.0.1:3306",
}

func init() {
	dbName := "test_createdb"

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)

	err := tp.CreateDatabase(dbName)

	if err != nil {
		log.Fatalln("Error creating database ", err)
	}
}

func TestCreateUser(t *testing.T) {
	dbName := "test_createdb"

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)

	err := tp.CreateUser(dbName, "mytestUser", "mytestPass")
	if err != nil {
		t.Errorf("Error creating user %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)

	err := tp.DeleteUser("mytestUser")
	if err != nil {
		t.Errorf("Error deleting user %v", err)
	}
}

func TestDeleteTheDatabase(t *testing.T) {
	dbName := "test_createdb"

	tp := New(mysqlConConfig.User, mysqlConConfig.Pass, mysqlConConfig.Host)

	err := tp.DeleteDatabase(dbName)
	if err != nil {
		t.Errorf("Error deleting database %v", err)
	}
}
