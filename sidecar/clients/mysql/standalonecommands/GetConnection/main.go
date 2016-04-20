package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/config"
	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/sidecar/clients/util"
	"github.com/pivotal-golang/lager"
)

func getDBNameFromId(id string) string {
	dbName := "d" + strings.Replace(id, "-", "", -1)
	dbName = strings.Replace(dbName, ";", "", -1)

	return dbName
}
func getMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

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
	Database string `json:"database"`
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

	username, err := getMD5Hash(os.Args[2])
	if err != nil {
		util.WriteError(err.Error(), 4, 500)
		os.Exit(4)
	}
	if len(username) > 16 {
		username = username[:16]
	}

	dbName := "d" + strings.Replace(os.Args[1], "-", "", -1)

	userCreated, err := provisioner.IsUserCreated(username)

	if err != nil {
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}

	if !userCreated {
		util.WriteError("Error 5: User not created", 5, 500)
		os.Exit(5)
	}

	mysqlConfigResp := MysqlConfigResp{Host: mysqlconfig.Host, Hostname: mysqlconfig.Host,
		Port: mysqlconfig.Port, User: username, Username: username,
		Database: dbName,
	}

	util.WriteSuccess(mysqlConfigResp, 200)
}
