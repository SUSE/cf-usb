package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"os"
	"regexp"
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

type WorkspaceOk struct {
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

	dbName := "d" + strings.Replace(os.Args[1], "-", "", -1)

	dbCreated, err := provisioner.IsDatabaseCreated(dbName)

	if err != nil {
		util.WriteError(err.Error(), 3, 500)
		os.Exit(3)
	}

	if !dbCreated {
		util.WriteError("Database not found", 4, 500)
		os.Exit(4)
	}
	util.WriteSuccess(WorkspaceOk{Message: "Database found"}, 200)
}
