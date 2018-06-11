package ccapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager"
)

//USBServiceBroker is the  interface to use for creating and relating with a a Service Broker
type USBServiceBroker interface {
	Create(name BrokerName, url, username, password string) error
	Delete(name BrokerName) error
	Update(serviceBrokerGUID BrokerGUID, name BrokerName, url, username, password string) error
	UpdateAll(url, username, password string) error
	EnableServiceAccess(ServiceGUID) error
	GetServiceBrokerGUIDByName(BrokerName) (BrokerGUID, error)
	GetServiceGUIDByName(ServiceName) (ServiceGUID, error)
	CheckServiceNameExists(ServiceName) (bool, error)
	CheckServiceInstancesExist(ServiceName) bool
}

//ServiceBroker is the definition of ServiceBroker type
type ServiceBroker struct {
	client         httpclient.HTTPClient
	tokenGenerator uaaapi.GetTokenInterface
	ccAPI          string
	logger         lager.Logger
}

// BrokerGUID is a specialized type for the GUID of a service broker
type BrokerGUID string

// BrokerName is a specialized type for the name of the service broker
type BrokerName string

//A BrokerEntity a broker-specific entity definition
type BrokerEntity struct {
	Name         BrokerName `json:"name,omitempty"`
	BrokerURL    string     `json:"broker_url"`
	AuthUsername string     `json:"auth_username"`
	AuthPassword string     `json:"auth_password"`
}

//BrokerResources holds the resources for the broker. Is mapped to json:resources
type BrokerResources struct {
	Resources []BrokerResource `json:"resources"`
}

//BrokerResource holds the broker metadata. Is mapped to json:metadata
type BrokerResource struct {
	Metadata struct {
		GUID BrokerGUID `json:"guid"`
	} `json:"metadata"`
	Entity BrokerEntity `json:"entity"`
}

//ServiceInstanceResources holds the service instance resources
type ServiceInstanceResources struct {
	Resources []ServiceInstance `json:"resources"`
}

//ServiceInstance holds the metadata and entity of service instance
type ServiceInstance struct {
	Metadata struct {
		GUID string `json:"guid"`
	} `json:"metadata"`
	Value struct {
		Name            ServiceInstanceName `json:"name"`
		ServicePlanGUID PlanGUID            `json:"service_plan_guid"`
	} `json:"entity"`
}

// A ServiceInstanceName is the name of a service instance
type ServiceInstanceName string

//NewServiceBroker creates and returns ServiceBroker
func NewServiceBroker(client httpclient.HTTPClient, token uaaapi.GetTokenInterface, ccAPI string, logger lager.Logger) USBServiceBroker {
	return &ServiceBroker{
		client:         client,
		tokenGenerator: token,
		ccAPI:          ccAPI,
		logger:         logger.Session("cc-service-broker-client", lager.Data{"cc-api": ccAPI}),
	}
}

//Create creates a service broker and returns an error if it fails
func (sb *ServiceBroker) Create(name BrokerName, url, username, password string) error {
	log := sb.logger.Session("create-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

	path := "/v2/service_brokers"
	body := &BrokerEntity{Name: name, BrokerURL: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Authorization": string(token),
	}

	log.Debug("preparing-request", lager.Data{"request-content": string(values)})

	request := httpclient.Request{Verb: "POST", Endpoint: sb.ccAPI, APIURL: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}

//Update updates a service broker
func (sb *ServiceBroker) Update(serviceBrokerGUID BrokerGUID, name BrokerName, url, username, password string) error {
	log := sb.logger.Session("update-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Authorization": string(token),
	}

	path := fmt.Sprintf("/v2/service_brokers/%s", serviceBrokerGUID)
	body := BrokerEntity{Name: name, BrokerURL: url, AuthUsername: username, AuthPassword: password}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	log.Debug("preparing-request", lager.Data{"request-content": string(values)})

	request := httpclient.Request{Verb: "PUT", Endpoint: sb.ccAPI, APIURL: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}

//UpdateAll updates all service brokers with the given url and username to the given password
func (sb *ServiceBroker) UpdateAll(url, username, password string) error {
	log := sb.logger.Session("update-all-brokers", lager.Data{"url": url})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Authorization": string(token),
	}

	path := "/v2/service_brokers"
	request := httpclient.Request{Verb: "GET", Endpoint: sb.ccAPI, APIURL: path, Headers: headers, StatusCode: 200}
	log.Info("list-brokers", lager.Data{"path": path})
	response, err := sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Debug("cc-response", lager.Data{"response": string(response)})

	resources := &BrokerResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return err
	}

	for _, resource := range resources.Resources {
		if resource.Entity.BrokerURL == url && resource.Entity.AuthUsername == username {
			err = sb.Update(resource.Metadata.GUID, resource.Entity.Name, url, username, password)
			if err != nil {
				return err
			}
		} else {
			log.Debug("skipping-broker", lager.Data{"broker": resource.Entity})
		}
	}

	log.Info("finished-cc-request")
	return nil
}

//EnableServiceAccess enables service access for the service having serviceName
func (sb *ServiceBroker) EnableServiceAccess(serviceGUID ServiceGUID) error {
	log := sb.logger.Session("enableservice-access", lager.Data{"service": serviceGUID})
	log.Debug("starting")

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccAPI, log)

	err := sp.Update(serviceGUID)
	if err != nil {
		return err
	}

	log.Debug("finished")

	return nil
}

//GetServiceBrokerGUIDByName obtains the broker guid corresponding to the passed name
func (sb *ServiceBroker) GetServiceBrokerGUIDByName(name BrokerName) (BrokerGUID, error) {
	log := sb.logger.Session("get-service-broker-guid-by-name", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/v2/service_brokers?q=name:%s", name)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Accept":       "application/json; charset=utf-8",
	}

	log.Debug("preparing-request", lager.Data{"path": path, "headers": headers})

	headers["Authorization"] = string(token)

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccAPI, APIURL: path, Headers: headers, StatusCode: 200}

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

	guid := resources.Resources[0].Metadata.GUID
	log.Debug("found", lager.Data{"service-broker-guid": guid})

	return guid, nil
}

// GetServiceGUIDByName returns the GUID of any services with the given name
func (sb *ServiceBroker) GetServiceGUIDByName(name ServiceName) (ServiceGUID, error) {
	log := sb.logger.Session("get-service-guid-by-name", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		log.Error("check-service-name-exists", err)
		return ServiceGUID(""), err
	}

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccAPI, log)
	guid, err := sp.GetServiceGUIDByLabel(name, token)
	if err != nil {
		log.Error("get-service-guid-by-label", err)
		return ServiceGUID(""), err
	}

	return guid, nil
}

//CheckServiceNameExists checks if a service with the passed name is already defined
func (sb *ServiceBroker) CheckServiceNameExists(name ServiceName) (bool, error) {
	guid, err := sb.GetServiceGUIDByName(name)
	if err != nil {
		return false, err
	}
	return (guid != ""), nil
}

//CheckServiceInstancesExist checks if a service instance with the passed name is already registered
func (sb *ServiceBroker) CheckServiceInstancesExist(serviceName ServiceName) bool {
	exist := false
	log := sb.logger.Session("check-service-instances-exist", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		log.Error("get-token-error", err)
		return false
	}

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccAPI, log)
	serviceGUID, err := sp.GetServiceGUIDByLabel(serviceName, token)
	if err != nil {
		log.Error("check-service-instance-exists-get-service-guid-by-label", err)
		return false
	}

	servicePlans, err := sp.GetServicePlans(serviceGUID, token)
	if err != nil {
		log.Error("check-service-instance-exists-get-service-plans", err)
		return false
	}

	headers := map[string]string{
		"Authorization": string(token),
		"Content-Type":  "application/x-www-form-urlencoded; charset=UTF-8",
		"Accept":        "application/json; charset=utf-8",
	}

	for _, plan := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s/service_instances", plan.Metadata.GUID)

		log.Debug("preparing-request-service_instances", lager.Data{"path": path, "headers": headers})

		findRequest := httpclient.Request{Verb: "GET", Endpoint: sb.ccAPI, APIURL: path, Headers: headers, StatusCode: 200}

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

//Delete deletes the service with the given name
func (sb *ServiceBroker) Delete(name BrokerName) error {
	log := sb.logger.Session("delete-broker", lager.Data{"name": name})
	log.Debug("starting")

	guid, err := sb.GetServiceBrokerGUIDByName(name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/v2/service_brokers/%s", guid)
	values := ""
	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Authorization": string(token),
	}

	log.Debug("preparing-request", lager.Data{"request-content": string(values)})

	request := httpclient.Request{Verb: "DELETE", Endpoint: sb.ccAPI, APIURL: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 204}

	log.Info("starting-cc-request", lager.Data{"path": path})

	_, err = sb.client.Request(request)
	if err != nil {
		return err
	}

	log.Info("finished-cc-request")

	return nil
}
