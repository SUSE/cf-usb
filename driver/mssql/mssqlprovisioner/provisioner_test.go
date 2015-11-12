package mssqlprovisioner

import (
	"database/sql"
	"fmt"
	"github.com/pivotal-golang/lager/lagertest"
	"os"
	"testing"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mssql-provisioner")

var mssqlConConfig = map[string]string{}

func init() {
	mssqlConConfig["server"] = os.Getenv("MSSQL_HOST")
	mssqlConConfig["port"] = os.Getenv("MSSQL_PORT")
	mssqlConConfig["user id"] = os.Getenv("MSSQL_USER")
	mssqlConConfig["password"] = os.Getenv("MSSQL_PASS")
}

func checkMssqlServer(t *testing.T) {
	if os.Getenv("MSSQL_HOST") == "" {
		t.Skip("Skipping test as not all env variables are set:'MSSQL_USER','MSSQL_PASS','MSSQL_HOST','MSSQL_PORT'")
	}
}

func TestCreateDatabaseDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.create-db"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}
	defer mssqlProv.Close()

	// Act
	err = mssqlProv.CreateDatabase(dbName)

	// Assert
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}
	defer sqlClient.Exec("drop database [" + dbName + "]")

	row := sqlClient.QueryRow("SELECT count(*) FROM sys.databases where name = ?", dbName)
	dbCount := 0
	row.Scan(&dbCount)
	if dbCount == 0 {
		t.Errorf("Database was not created")
	}
}

func TestDeleteDatabaseDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.delete-db"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Database init error, %v", err)
	}
	defer mssqlProv.Close()

	err = mssqlProv.CreateDatabase(dbName)

	// Act

	err = mssqlProv.DeleteDatabase(dbName)

	// Assert
	if err != nil {
		t.Errorf("Database delete error, %v", err)
	}

	row := sqlClient.QueryRow("SELECT count(*) FROM sys.databases where name = ?", dbName)
	dbCount := 0
	row.Scan(&dbCount)
	if dbCount != 0 {
		t.Errorf("Database %s was not deleted", dbName)
	}
}

func TestCreateUserDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.create-db"
	userNanme := "cf-broker-testing.create-user"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)

	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}

	err = mssqlProv.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}

	// Act
	err = mssqlProv.CreateUser(dbName, userNanme, "passwordAa_0")

	// Assert
	if err != nil {
		t.Errorf("User create error, %v", err)
	}

	defer sqlClient.Exec("drop database [" + dbName + "]")

	row := sqlClient.QueryRow(fmt.Sprintf("select count(*)  from [%s].sys.database_principals  where name = ?", dbName), userNanme)
	dbCount := 0
	row.Scan(&dbCount)
	if dbCount == 0 {
		t.Errorf("User was not created")
	}
}

func TestDeleteUserDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.create-db"
	userNanme := "cf-broker-testing.create-user"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}
	defer mssqlProv.Close()

	err = mssqlProv.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}
	err = mssqlProv.CreateUser(dbName, userNanme, "passwordAa_0")
	if err != nil {
		t.Errorf("User create error, %v", err)
	}

	// Act
	exists, err := mssqlProv.IsUserCreated(dbName, userNanme)

	// Assert
	if err != nil {
		t.Errorf("IsUserCreated error, %v", err)
	}
	if !exists {
		t.Errorf("IsUserCreated returned false, expected true")
	}

	// Act
	err = mssqlProv.DeleteUser(dbName, userNanme)

	// Assert
	if err != nil {
		t.Errorf("User delete error, %v", err)
	}

	// Act
	exists, err = mssqlProv.IsUserCreated(dbName, userNanme)

	// Assert
	if err != nil {
		t.Errorf("IsUserCreated error, %v", err)
	}
	if exists {
		t.Errorf("IsUserCreated returned true, expected false")
	}

	defer sqlClient.Exec("drop database [" + dbName + "]")

	row := sqlClient.QueryRow(fmt.Sprintf("select count(*)  from [%s].sys.database_principals  where name = ?", dbName), userNanme)
	dbCount := 0
	row.Scan(&dbCount)
	if dbCount != 0 {
		t.Errorf("User was not deleted")
	}
}

func TestIsDatabaseCreatedDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.nonexisting-db"

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err := mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}

	// Act
	exists, err := mssqlProv.IsDatabaseCreated(dbName)

	// Assert
	if err != nil {
		t.Errorf("Check for database error, %v", err)
	}
	if exists {
		t.Errorf("Check for database error, expected false, but received true")
	}
}

func TestIsDatabaseCreatedDriver2(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.create-db"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}
	err = mssqlProv.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}

	// Act
	exists, err := mssqlProv.IsDatabaseCreated(dbName)

	// Assert
	if err != nil {
		t.Errorf("Check for database error, %v", err)
	}
	if !exists {
		t.Errorf("Check for database error, expected true, but received false")
	}

	defer sqlClient.Exec("drop database [" + dbName + "]")
}

func TestStressDriver(t *testing.T) {
	checkMssqlServer(t)

	dbName := "cf-broker-testing.create-db"
	dbNameA := "cf-broker-testing.create-db-A"
	dbName2 := "cf-broker-testing.create-db-2"
	userNanme := "cf-broker-testing.create-user"

	sqlClient, err := sql.Open("mssql", buildConnectionString(mssqlConConfig))
	defer sqlClient.Close()

	sqlClient.Exec("drop database [" + dbName + "]")
	sqlClient.Exec("drop database [" + dbNameA + "]")
	sqlClient.Exec("drop database [" + dbName2 + "]")

	logger = lagertest.NewTestLogger("process-controller")
	mssqlProv := NewMssqlProvisioner(logger)
	err = mssqlProv.Connect("mssql", mssqlConConfig)
	if err != nil {
		t.Errorf("Provisioner init error, %v", err)
	}

	err = mssqlProv.CreateDatabase(dbName)
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}

	err = mssqlProv.CreateDatabase(dbNameA)
	if err != nil {
		t.Errorf("Database create error, %v", err)
	}

	wait := make(chan bool)

	go func() {
		for i := 1; i < 8; i++ {

			err := mssqlProv.CreateDatabase(dbName2)
			if err != nil {
				t.Errorf("Database create error, %v", err)
				break
			}

			err = mssqlProv.DeleteDatabase(dbName2)
			if err != nil {
				t.Errorf("Database delete error, %v", err)
				break
			}
		}

		wait <- true
	}()

	go func() {
		for i := 1; i < 32; i++ {
			err = mssqlProv.CreateUser(dbName, userNanme, "passwordAa_0")
			if err != nil {
				t.Errorf("User create error, %v", err)
				break
			}

			err = mssqlProv.DeleteUser(dbName, userNanme)
			if err != nil {
				t.Errorf("User delete error, %v", err)
				break
			}

		}

		wait <- true
	}()

	go func() {
		for i := 1; i < 32; i++ {
			err = mssqlProv.CreateUser(dbNameA, userNanme, "passwordAa_0")
			if err != nil {
				t.Errorf("User create error, %v", err)
				break
			}

			err = mssqlProv.DeleteUser(dbNameA, userNanme)
			if err != nil {
				t.Errorf("User delete error, %v", err)
				break
			}

		}

		wait <- true
	}()

	<-wait
	<-wait
	<-wait

	sqlClient.Exec("drop database [" + dbName + "]")
	sqlClient.Exec("drop database [" + dbName2 + "]")
	sqlClient.Exec("drop database [" + dbNameA + "]")
}
