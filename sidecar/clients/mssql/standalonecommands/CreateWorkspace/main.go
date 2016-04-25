package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hpcloud/cf-usb/sidecar/clients/mssql/mssqlprovisioner"
	"github.com/hpcloud/cf-usb/sidecar/clients/util"
	"github.com/pivotal-golang/lager"
)

func getDBNameFromId(id string) string {
	dbName := "d" + strings.Replace(id, "-", "", -1)
	dbName = strings.Replace(dbName, ";", "", -1)

	return dbName
}

type DbCreated struct {
	Message string `json:"message"`
}

func main() {
	if len(os.Args) < 2 {
		util.WriteError("No database name provided", 1, 500)
		os.Exit(1)

	}
	mhost := os.Getenv("MSSQL_HOST")
	mport := os.Getenv("MSSQL_PORT")
	muser := os.Getenv("MSSQL_USER")
	mpass := os.Getenv("MSSQL_PASS")

	if mhost == "" || mport == "" || muser == "" || mpass == "" {
		util.WriteError("MSSQL_HOST, MSSQL_PORT, MSSQL_USER and MSSQL_PASS env vars not set!", 2, 500)
		os.Exit(2)
	}

	port, err := strconv.Atoi(os.Getenv("MSSQL_PORT"))
	if err != nil {
		os.Exit(1)
	}

	var mssqlConConfig = map[string]string{}
	mssqlConConfig["server"] = mhost
	mssqlConConfig["port"] = strconv.Itoa(port)
	mssqlConConfig["user id"] = muser
	mssqlConConfig["password"] = mpass

	loger := lager.NewLogger("stdout")

	provisioner := mssqlprovisioner.NewMssqlProvisioner(loger)
	provisioner.Connect("mssql", mssqlConConfig)

	err = provisioner.CreateDatabase(getDBNameFromId(os.Args[1]))

	if err != nil {
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}
	dbCreated := DbCreated{Message: fmt.Sprintf("Created db: %s", getDBNameFromId(os.Args[1]))}
	util.WriteSuccess(dbCreated, 200)
	//if everything is ok exit status will be 0
}
