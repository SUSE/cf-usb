package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/config"
	"github.com/hpcloud/cf-usb/sidecar/clients/mysql/mysqlprovisioner"
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

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}

type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

//this with ErrorResp will transform an error in a json for server use
func WriteError(err error) {
	var errResp = ErrorResp{}

	if strings.HasPrefix(err.Error(), "Error") && strings.Contains(err.Error(), ":") {
		rez, errConv := strconv.Atoi(strings.TrimPrefix(strings.Split(err.Error(), ":")[0], "Error "))
		if errConv != nil {
			errResp.Code = 500
			errResp.Message = err.Error()
		} else {
			errResp.Code = rez
			errResp.Message = strings.Split(err.Error(), ":")[1]
		}
	} else {
		errResp.Code = 500
		errResp.Message = err.Error()
	}
	strResp, _ := json.Marshal(errResp)
	os.Stdout.WriteString("500\r\n") //error code
	os.Stdout.WriteString(string(strResp))
}

func main() {
	if len(os.Args) < 3 {
		WriteError(errors.New(fmt.Sprintf("No database name or username provided -d %d\n", len(os.Args))))
		os.Exit(1)
	}
	mhost := os.Getenv("MYSQL_HOST")
	mport := os.Getenv("MYSQL_PORT")
	muser := os.Getenv("MYSQL_USER")
	mpass := os.Getenv("MYSQL_PASS")

	if mhost == "" || mport == "" || muser == "" || mpass == "" {
		WriteError(errors.New("MYSQL_HOST, MYSQL_PORT, MYSQL_USER and MYSQL_PASS env vars not set!"))
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
		WriteError(err)
		os.Exit(4)
	}
	if len(username) > 16 {
		username = username[:16]
	}

	dbName := "d" + strings.Replace(os.Args[1], "-", "", -1)

	userCreated, err := provisioner.IsUserCreated(username)

	if err != nil {
		WriteError(err)
		os.Exit(3)
	}

	if !userCreated {
		WriteError(errors.New("Error 5: User not created"))
		os.Exit(5)
	}

	mysqlConfigResp := MysqlConfigResp{Host: mysqlconfig.Host, Hostname: mysqlconfig.Host,
		Port: mysqlconfig.Port, User: username, Username: username,
		Database: dbName,
	}
	bRez, err := json.Marshal(mysqlConfigResp)
	if err != nil {
		WriteError(errors.New("Error 500: Unknown error"))
		os.Exit(4)
	}
	fmt.Println(string(bRez))
}
