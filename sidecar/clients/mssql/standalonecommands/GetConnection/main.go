package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/cf-usb/sidecar/clients/mssql/config"

	"github.com/hpcloud/cf-usb/sidecar/clients/mssql/mssqlprovisioner"
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

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}

func main() {
	if len(os.Args) < 3 {
		util.WriteError(fmt.Sprintf("No database name or username provided -d %d\n", len(os.Args)), 1, 500)
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

	dbName := "d" + strings.Replace(os.Args[1], "-", "", -1)

	username, err := getMD5Hash(os.Args[2])
	if err != nil {
		util.WriteError(err.Error(), 4, 500)
		os.Exit(4)
	}
	if len(username) > 16 {
		username = username[:16]
	}

	userCreated, err := provisioner.IsUserCreated(dbName, username)

	if err != nil {
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}

	if !userCreated {
		util.WriteError("Error 5: User not created", 5, 500)
		os.Exit(5)
	}

	mssqlResp := config.MssqlBindingCredentials{
		Host:     mhost,
		Port:     port,
		Username: username,
	}

	util.WriteSuccess(mssqlResp, 200)
}
