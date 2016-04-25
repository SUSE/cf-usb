package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/config"
	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/sidecar/clients/util"
	"github.com/pivotal-golang/lager"
)

type MysqlDriver struct {
	logger lager.Logger
	conf   config.MysqlDriverConfig
	db     mysqlprovisioner.MysqlProvisionerInterface
}

type MysqlConfigResp struct {
	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func WriteInternalError() {
	os.Stdout.WriteString("{\"code\":500,\"message\":\"unknown error\" ")
}

func main() {
	if len(os.Args) < 3 {
		util.WriteError(fmt.Sprintf("No database name or username provided -d %d\n", len(os.Args)), 1, 500)
		os.Exit(1)

	}
	mhost := os.Getenv("MYSQL_HOST")
	mport := os.Getenv("MYSQL_PORT")
	muser := os.Getenv("MYSQL_USER")
	mpass := os.Getenv("MYSQL_PASS")

	if mhost == "" || mport == "" || muser == "" || mpass == "" {
		util.WriteError("MYSQL_HOST, MYSQL_PORT, MYSQL_USER and MYSQL_PASS env vars not set!", 2, 417) //expectation failed
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
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}
	if len(username) > 16 {
		username = username[:16]
	}
	password, _ := util.SecureRandomString(32)
	dbName := "d" + strings.Replace(os.Args[1], "-", "", -1)

	err = provisioner.CreateUser(dbName, username, password)

	if err != nil {
		util.WriteError(err.Error(), 4, 410)
		os.Exit(4)
	}

	mysqlConfigResp := MysqlConfigResp{Hostname: mysqlconfig.Host, Host: mysqlconfig.Host,
		Port: mysqlconfig.Port, Username: username, User: username,
		Password: password, Database: dbName}
	util.WriteSuccess(mysqlConfigResp, 200)

}
