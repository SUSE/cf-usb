package uaaapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager"
)

// A BearerToken is a string that can be set in a HTTP header as proof of the token
type BearerToken string

//GetTokenInterface is the interface used to obtain token from uaa api
type GetTokenInterface interface {
	GetToken() (BearerToken, error)
}

//Token defines auth token basic struct
type Token struct {
	AccessToken string `json:"access_token"`
	ExpireTime  int    `json:"expires_in"`
}

//Generator defines the generation of tokens
type Generator struct {
	tokenURL     string
	clientID     string
	clientSecret string
	client       httpclient.HTTPClient
	logger       lager.Logger
}

//NewTokenGenerator creates and returns a TokenGenerator
func NewTokenGenerator(tokenURL, clientID, clientSecret string, client httpclient.HTTPClient, logger lager.Logger) GetTokenInterface {
	return &Generator{
		tokenURL:     tokenURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       client,
		logger:       logger.Session("uaa-token-generator"),
	}
}

//GetToken obtains the token from the generator
func (generator *Generator) GetToken() (BearerToken, error) {
	log := generator.logger.Session("fetch-token", lager.Data{"uaa-api": generator.tokenURL})
	log.Debug("starting")
	defer log.Debug("finished")

	valuesBody := url.Values{}
	valuesBody.Add("grant_type", "client_credentials")
	requestBody := valuesBody.Encode()

	tokenURL := fmt.Sprintf("/oauth/token")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	credentials := httpclient.BasicAuth{Username: generator.clientID, Password: generator.clientSecret}
	request := httpclient.Request{Verb: "POST", Endpoint: generator.tokenURL, APIURL: tokenURL, Body: strings.NewReader(requestBody), Headers: headers, Credentials: &credentials, StatusCode: 200}

	log.Info("starting-uaa-request", lager.Data{"path": tokenURL, "verb": "GET"})

	response, err := generator.client.Request(request)
	if err != nil {
		return "", err
	}

	log.Info("finished-uaa-request")

	token := &Token{}
	err = json.Unmarshal(response, token)
	if err != nil {
		return "", err
	}

	bearerToken := fmt.Sprintf("bearer %v", token.AccessToken)
	return BearerToken(bearerToken), nil
}
