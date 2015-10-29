package cc_api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaa_api"
	"github.com/pivotal-golang/lager"
)

type ServiceBrokerInterface interface {
	Create(name, url, username, password string) error
	Update(serviceBrokerGuid, name, url, username, password string) error
}

type ServiceBroker struct {
	client         httpclient.HttpClient
	tokenGenerator uaa_api.GetTokenInterface
	ccApi          string
	logger         lager.Logger
}

type BrokerValues struct {
	Name         string `json:"name"`
	BrokerUrl    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
}

func NewServiceBroker(client httpclient.HttpClient, token uaa_api.GetTokenInterface, ccApi string, logger lager.Logger) ServiceBrokerInterface {
	return &ServiceBroker{
		client:         client,
		tokenGenerator: token,
		ccApi:          ccApi,
		logger:         logger,
	}
}

func (sb *ServiceBroker) Create(name, url, username, password string) error {
	path := "/v2/service_brokers"
	body := &BrokerValues{Name: name, BrokerUrl: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sb.logger.Debug("create-service-broker", lager.Data{"service broker name": name, "content": string(values)})

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	sb.logger.Debug("create-service-broker", lager.Data{"token recieved": token})

	headers := make(map[string]string)
	headers["Authorization"] = token

	request := httpclient.Request{Verb: "POST", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	return nil
}

func (sb *ServiceBroker) Update(serviceBrokerGuid, name, url, username, password string) error {
	path := fmt.Sprintf("/v2/service_brokers/%s", serviceBrokerGuid)
	body := BrokerValues{Name: name, BrokerUrl: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sb.logger.Debug("update-service-broker", lager.Data{"service broker name": name})

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	request := httpclient.Request{Verb: "PUT", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 200}

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	return nil
}
