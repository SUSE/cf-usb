package ccapi

import (
	"encoding/json"
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
	GetServices() (*Services, error)
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

type Services struct {
	Resources []ServiceInfo `json:"resources"`
}

type ServiceInfo struct {
	Metadata ServiceMetadata `json:"metadata"`
	Entity   Service         `json:"entity"`
}

type ServiceMetadata struct {
	Guid       string `json:"guid"`
	Url        string `json:"url,omitempty"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type Service struct {
	Name        string `json:"label"`
	Description string `json:"description"`
	Unique_Id   string `json:"unique_id,omitempty"`
}

func NewServiceBroker(client httpclient.HttpClient, token uaaapi.GetTokenInterface, ccApi string, logger lager.Logger) ServiceBrokerInterface {
	return &ServiceBroker{
		client:         client,
		tokenGenerator: token,
		ccApi:          ccApi,
		logger:         logger.Session("cc-service-broker-client", lager.Data{"cc-api": ccApi}),
	}
}

func (sb *ServiceBroker) Create(name, url, username, password string) error {
	log := sb.logger.Session("create-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

	path := "/v2/service_brokers"
	body := &BrokerValues{Name: name, BrokerUrl: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	log.Debug("preparing-request", lager.Data{"request-content": string(values), "headers": headers})

	request := httpclient.Request{Verb: "POST", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}

func (sb *ServiceBroker) Update(serviceBrokerGuid, name, url, username, password string) error {
	log := sb.logger.Session("update-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

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

	log.Debug("preparing-request", lager.Data{"request-content": string(values), "headers": headers})

	request := httpclient.Request{Verb: "PUT", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}

func (sb *ServiceBroker) EnableServiceAccess(serviceId string) error {
	log := sb.logger.Session("enableservice-access", lager.Data{"serviceID": serviceId})
	log.Debug("starting")

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccApi, log)

	err := sp.Update(serviceId)
	if err != nil {
		return err
	}

	log.Debug("finished")

	return nil
}

func (sb *ServiceBroker) GetServiceBrokerGuidByName(name string) (string, error) {
	log := sb.logger.Session("get-service-broker-guid-by-name", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/v2/service_brokers?q=name:%s", name)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	log.Debug("preparing-request", lager.Data{"path": path, "headers": headers})

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path})

	response, err := sb.client.Request(findRequest)
	if err != nil {
		return "", err
	}

	log.Debug("cc-reponse", lager.Data{"response": string(response)})
	log.Info("finished-cc-request")

	resources := &BrokerResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return "", err
	}

	if len(resources.Resources) == 0 {
		log.Debug("not-found")
		return "", nil
	}

	guid := resources.Resources[0].Values.Guid
	log.Debug("found", lager.Data{"service-broker-guid": guid})

	return guid, nil
}

func (sb *ServiceBroker) GetServices() (*Services, error) {
	log := sb.logger.Session("get-services")
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return nil, err
	}

	path := "/v2/services"

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	log.Debug("preparing-request", lager.Data{"path": path, "headers": headers})
	findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path})

	response, err := sb.client.Request(findRequest)
	if err != nil {
		return nil, err
	}

	log.Debug("cc-reponse", lager.Data{"response": string(response)})
	log.Info("finished-cc-request")

	resources := &Services{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
