package ccapi

import (
	"encoding/json"
	"fmt"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager"
)

//GetInfoInterface is the interface for providing information about token endpoint
type GetInfoInterface interface {
	GetTokenEndpoint() (string, error)
}

//GetInfo is the definition of the GetInfo type
type GetInfo struct {
	ccAPI  string
	client httpclient.HTTPClient
	logger lager.Logger
}

//GetInfoResponse is the structure of token endpoint
type GetInfoResponse struct {
	TokenEndpoint string `json:"token_endpoint"`
}

//NewGetInfo instantiates a new GetInfo
func NewGetInfo(ccAPI string, client httpclient.HTTPClient, logger lager.Logger) GetInfoInterface {
	return &GetInfo{
		ccAPI:  ccAPI,
		client: client,
		logger: logger.Session("cc-info"),
	}
}

//GetTokenEndpoint obtains the endpoint from GetInfo
func (info *GetInfo) GetTokenEndpoint() (string, error) {
	log := info.logger.Session("get-token-endpoint", lager.Data{"cc-api": info.ccAPI})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/info")

	request := httpclient.Request{Verb: "GET", Endpoint: info.ccAPI, APIURL: path, StatusCode: 200}

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

	if data.TokenEndpoint == "" {
		return "", fmt.Errorf("UAA token endpoint missing")
	}
	return data.TokenEndpoint, nil
}
