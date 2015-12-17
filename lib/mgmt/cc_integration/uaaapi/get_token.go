package uaaapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager"
)

type GetTokenInterface interface {
	GetToken() (string, error)
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpireTime  int    `json:"expires_in"`
}

type Generator struct {
	tokenUrl     string
	clientId     string
	clientSecret string
	client       httpclient.HttpClient
	logger       lager.Logger
}

func NewTokenGenerator(tokenUrl, clientId, clientSecret string, client httpclient.HttpClient, logger lager.Logger) GetTokenInterface {
	return &Generator{
		tokenUrl:     tokenUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		client:       client,
		logger:       logger.Session("uaa-token-generator"),
	}
}

func (generator *Generator) GetToken() (string, error) {
	log := generator.logger.Session("fetch-token", lager.Data{"uaa-api": generator.tokenUrl})
	log.Debug("starting")

	valuesBody := url.Values{}
	valuesBody.Add("grant_type", "client_credentials")
	requestBody := valuesBody.Encode()

	tokenURL := fmt.Sprintf("/oauth/token")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	credentials := httpclient.BasicAuth{Username: generator.clientId, Password: generator.clientSecret}
	request := httpclient.Request{Verb: "POST", Endpoint: generator.tokenUrl, ApiUrl: tokenURL, Body: strings.NewReader(requestBody), Headers: headers, Credentials: &credentials, StatusCode: 200}

	log.Info("starting-uaa-request", lager.Data{"path": tokenURL, "verb": "GET"})

	response, err := generator.client.Request(request)
	if err != nil {
		return "", err
	}

	log.Debug("uaa-reponse", lager.Data{"response": string(response)})
	log.Info("finished-uaa-request")

	token := &Token{}
	err = json.Unmarshal(response, token)
	if err != nil {
		return "", err
	}

	bearerToken := fmt.Sprintf("bearer %v", token.AccessToken)
	return bearerToken, nil
}
