package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type JsonResp struct {
	HttpCode int         `json:"httpCode"`
	Payload  interface{} `json:"payload"`
}

type ErrorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(errStr string, code int, httpCode int) {
	errMsg := ErrorMsg{}

	errMsg.Code = code
	errMsg.Message = errStr

	errResp := JsonResp{}
	errResp.HttpCode = httpCode
	errResp.Payload = errMsg
	strResp, _ := json.Marshal(errResp)
	os.Stdout.WriteString(string(strResp))
}

func WriteSuccess(payload interface{}, httpCode int) {

	jsonResp := JsonResp{}
	jsonResp.HttpCode = httpCode
	jsonResp.Payload = payload

	bytesOk, err := json.Marshal(jsonResp)
	if err != nil {
		WriteError("Unknown error", 5, 500)
		os.Exit(5)
	}

	os.Stdout.WriteString(string(bytesOk))
}

func GetDBNameFromId(id string) string {
	dbName := "d" + strings.Replace(id, "-", "", -1)
	dbName = strings.Replace(dbName, "`", "", -1)
	dbName = strings.Replace(dbName, ";", "", -1)
	if len(dbName) > 64 {
		dbName = dbName[:64]
	}

	return dbName
}
func GetMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

func SecureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}
