package cc_api

import (
	"encoding/json"
	"fmt"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager"
)

type GetInfoInterface interface {
	GetTokenEndpoint() (string, error)
}

type GetInfo struct {
	ccApi  string
	client httpclient.HttpClient
	logger lager.Logger
}

type GetInfoResponse struct {
	Token_endpoint string `json:"token_endpoint"`
}

func NewGetInfo(ccApi string, client httpclient.HttpClient, logger lager.Logger) GetInfoInterface {
	return &GetInfo{
		ccApi:  ccApi,
		client: client,
		logger: logger,
	}
}

func (info *GetInfo) GetTokenEndpoint() (string, error) {
	path := fmt.Sprintf("/v2/info")

	request := httpclient.Request{Verb: "GET", Endpoint: info.ccApi, ApiUrl: path, StatusCode: 200}

	response, err := info.client.Request(request)
	if err != nil {
		return "", err
	}

	var data GetInfoResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return "", err
	}

	info.logger.Info("get-info", lager.Data{"get token endpoint": data.Token_endpoint})

	return data.Token_endpoint, nil
}
