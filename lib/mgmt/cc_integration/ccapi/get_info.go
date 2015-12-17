package ccapi

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
	TokenEndpoint string `json:"token_endpoint"`
}

func NewGetInfo(ccApi string, client httpclient.HttpClient, logger lager.Logger) GetInfoInterface {
	return &GetInfo{
		ccApi:  ccApi,
		client: client,
		logger: logger.Session("cc-info"),
	}
}

func (info *GetInfo) GetTokenEndpoint() (string, error) {
	log := info.logger.Session("get-token-endpoint", lager.Data{"cc-api": info.ccApi})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/info")

	request := httpclient.Request{Verb: "GET", Endpoint: info.ccApi, ApiUrl: path, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path, "verb": "GET"})

	response, err := info.client.Request(request)
	if err != nil {
		return "", err
	}

	log.Debug("cc-reponse", lager.Data{"response": string(response)})
	log.Info("finished-cc-request")

	var data GetInfoResponse
	err = json.Unmarshal(response, &data)
	if err != nil {
		return "", err
	}

	return data.TokenEndpoint, nil
}
