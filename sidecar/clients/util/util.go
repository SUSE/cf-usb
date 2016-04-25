package util

import (
	"encoding/json"
	"os"
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
