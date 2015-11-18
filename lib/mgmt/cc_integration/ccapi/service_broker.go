package ccapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager"
)

type ServiceBrokerInterface interface {
	Create(name, url, username, password string) error
	Update(serviceBrokerGuid, name, url, username, password string) error
	EnableServiceAccess(serviceId string) error
	GetServiceBrokerGuidByName(name string) (string, error)
}

type ServiceBroker struct {
	client         httpclient.HttpClient
	tokenGenerator uaaapi.GetTokenInterface
	ccApi          string
	logger         lager.Logger
}

type BrokerValues struct {
	Name         string `json:"name"`
	BrokerUrl    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
}

type BrokerResources struct {
	Resources []BrokerResource `json:"resources"`
}

type BrokerResource struct {
	Values BrokerMetadata `json:"metadata"`
}

type BrokerMetadata struct {
	Guid string `json:"guid"`
}

func NewServiceBroker(client httpclient.HttpClient, token uaaapi.GetTokenInterface, ccApi string, logger lager.Logger) ServiceBrokerInterface {
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
	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	path := fmt.Sprintf("/v2/service_brokers/%s", serviceBrokerGuid)
	body := BrokerValues{Name: name, BrokerUrl: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sb.logger.Debug("update-service-broker", lager.Data{"service broker name": name})

	request := httpclient.Request{Verb: "PUT", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 200}

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	return nil
}

func (sb *ServiceBroker) EnableServiceAccess(serviceId string) error {
	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccApi, sb.logger)

	err := sp.Update(serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (sb *ServiceBroker) GetServiceBrokerGuidByName(name string) (string, error) {
	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/v2/service_brokers?q=name:%s", name)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

	response, err := sb.client.Request(findRequest)
	if err != nil {
		return "", err
	}

	resources := &BrokerResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return "", err
	}

	if len(resources.Resources) == 0 {
		return "", errors.New(fmt.Sprintf("Service broker %s not found", name))
	}

	guid := resources.Resources[0].Values.Guid

	sb.logger.Debug("get-service-broker", lager.Data{"service broker guid": guid})

	return guid, nil
}
