package main

import (
	"fmt"
	"os"

	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/config"
	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/sidecar/clients/util"
	"github.com/pivotal-golang/lager"
)

type DeletedOK struct {
	Message string `json:"message"`
}

type MysqlDriver struct {
	logger lager.Logger
	conf   config.MysqlDriverConfig
	db     mysqlprovisioner.MysqlProvisionerInterface
}

func main() {
	if len(os.Args) < 2 {
		util.WriteError(fmt.Sprintf("No username provided -d %d", len(os.Args)), 1, 500)
		os.Exit(1)

	}
	mhost := os.Getenv("MYSQL_HOST")
	mport := os.Getenv("MYSQL_PORT")
	muser := os.Getenv("MYSQL_USER")
	mpass := os.Getenv("MYSQL_PASS")

	if mhost == "" || mport == "" || muser == "" || mpass == "" {
		util.WriteError("MYSQL_HOST, MYSQL_PORT, MYSQL_USER and MYSQL_PASS env vars not set!", 2, 500)
		os.Exit(2)
	}

	mysqlconfig := config.MysqlDriverConfig{}
	mysqlconfig.Host = os.Getenv("MYSQL_HOST")
	mysqlconfig.Pass = os.Getenv("MYSQL_PASS")
	mysqlconfig.Port = os.Getenv("MYSQL_PORT")
	mysqlconfig.User = os.Getenv("MYSQL_USER")

	loger := lager.NewLogger("stdout")

	provisioner := mysqlprovisioner.New(loger)
	provisioner.Connect(mysqlconfig)

	username, err := util.GetMD5Hash(os.Args[2])
	if err != nil {
		util.WriteError(err.Error(), 4, 500)
		os.Exit(4)
	}
	if len(username) > 16 {
		username = username[:16]
	}

	err = provisioner.DeleteUser(username)

	if err != nil {
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}

	util.WriteSuccess(DeletedOK{Message: "User deleted"}, 200)

	//if everything is ok exit status will be 0
}
