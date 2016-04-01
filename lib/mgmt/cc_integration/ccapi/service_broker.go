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
	Delete(name string) error
	Update(serviceBrokerGuid, name, url, username, password string) error
	EnableServiceAccess(serviceId string) error
	GetServiceBrokerGuidByName(name string) (string, error)
	CheckServiceNameExists(name string) bool
	CheckServiceInstancesExist(serviceName string) bool
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

type ServiceInstanceResources struct {
	Resources []ServiceInstance `json:"resources"`
}

type ServiceInstance struct {
	Meta  ServiceInstanceMetadata `json:"metadata"`
	Value ServiceInstanceEntity   `json:"entity"`
}

type ServiceInstanceMetadata struct {
	Guid string `json:"guid"`
}

type ServiceInstanceEntity struct {
	Name            string `json:"name"`
	ServicePlanGuid string `json:"service_plan_guid"`
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

func (sb *ServiceBroker) EnableServiceAccess(serviceName string) error {
	log := sb.logger.Session("enableservice-access", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccApi, log)

	err := sp.Update(serviceName)
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

func (sb *ServiceBroker) CheckServiceNameExists(name string) bool {
	exist := false
	log := sb.logger.Session("check-service-name-exists", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		log.Error("check-service-name-exists", err)
		return false
	}

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccApi, log)
	guid, err := sp.GetServiceGuidByLabel(name, token)
	if err != nil {
		log.Error("get-service-guid-by-label", err)
	}
	if guid != "" {
		exist = true
	}
	log.Debug(fmt.Sprintf("check service name %s exists complete - returning %t", name, exist))

	return exist
}

func (sb *ServiceBroker) CheckServiceInstancesExist(serviceName string) bool {
	exist := false
	log := sb.logger.Session("check-service-instances-exist", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		log.Error("get-token-error", err)
		return false
	}

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccApi, log)
	serviceGuid, err := sp.GetServiceGuidByLabel(serviceName, token)
	if err != nil {
		log.Error("check-service-instance-exists-get-service-guid-by-label", err)
		return false
	}

	servicePlans, err := sp.GetServicePlans(serviceGuid, token)
	if err != nil {
		log.Error("check-service-instance-exists-get-service-plans", err)
		return false
	}
	
	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	for _, plan := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s/service_instances", plan.Values.Guid)

		log.Debug("preparing-request-service_instances", lager.Data{"path": path, "headers": headers})

		findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

		log.Info("starting-cc-request-service_instances", lager.Data{"path": path})

		responseInstances, err := sb.client.Request(findRequest)
		if err != nil {
			log.Error("client-request-error-service_instances", err)
			return false
		}

		resourcesInstances := &ServiceInstanceResources{}
		err = json.Unmarshal(responseInstances, &resourcesInstances)
		if err != nil {
			log.Error("unmarshal-service-instances-resources", err)
			return false
		}
		if len(resourcesInstances.Resources) > 0 {
			exist = true
			break
		}
	}

	return exist
}

func (sb *ServiceBroker) Delete(name string) error {
	log := sb.logger.Session("delete-broker", lager.Data{"name": name})
	log.Debug("starting")

	guid, err := sb.GetServiceBrokerGuidByName(name)
	if err != nil {
		return err
	}

	path := "/v2/service_brokers/" + guid
	values := ""
	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	log.Debug("preparing-request", lager.Data{"request-content": string(values), "headers": headers})

	request := httpclient.Request{Verb: "DELETE", Endpoint: sb.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 204}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}
