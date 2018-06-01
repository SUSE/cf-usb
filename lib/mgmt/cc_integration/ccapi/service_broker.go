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
	Create(name, url, username, password string) error
	Delete(name string) error
	Update(serviceBrokerGUID, name, url, username, password string) error
	UpdateAll(url, username, password string) error
	EnableServiceAccess(serviceID string) error
	GetServiceBrokerGUIDByName(name string) (string, error)
	CheckServiceNameExists(name string) (bool, error)
	CheckServiceInstancesExist(serviceName string) bool
}

//ServiceBroker is the definition of ServiceBroker type
type ServiceBroker struct {
	client         httpclient.HTTPClient
	tokenGenerator uaaapi.GetTokenInterface
	ccAPI          string
	logger         lager.Logger
}

//BrokerValues is the type defining BrokerValues and maped to json values for BrokerValues
type BrokerValues struct {
	Name         string `json:"name,omitempty"`
	BrokerURL    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
}

//BrokerResources holds the resources for the broker. Is mapped to json:resources
type BrokerResources struct {
	Resources []BrokerResource `json:"resources"`
}

//BrokerResource holds the broker metadata. Is mapped to json:metadata
type BrokerResource struct {
	Meta  BrokerMetadata `json:"metadata"`
	Value BrokerValues   `json:"entity"`
}

//BrokerMetadata mapped to json:guid
type BrokerMetadata struct {
	GUID string `json:"guid"`
}

//ServiceInstanceResources holds the service instance resources
type ServiceInstanceResources struct {
	Resources []ServiceInstance `json:"resources"`
}

//ServiceInstance holds the metadata and entity of service instance
type ServiceInstance struct {
	Meta  ServiceInstanceMetadata `json:"metadata"`
	Value ServiceInstanceEntity   `json:"entity"`
}

//ServiceInstanceMetadata hold the metadata and is mapped to json:guid
type ServiceInstanceMetadata struct {
	GUID string `json:"guid"`
}

//ServiceInstanceEntity holds the Name and ServicePlanGUID for service instance. Is mapped to json
type ServiceInstanceEntity struct {
	Name            string `json:"name"`
	ServicePlanGUID string `json:"service_plan_guid"`
}

//NewServiceBroker creates and returns ServiceBroker
func NewServiceBroker(client httpclient.HTTPClient, token uaaapi.GetTokenInterface, ccAPI string, logger lager.Logger) USBServiceBroker {
	return &ServiceBroker{
		client:         client,
		tokenGenerator: token,
		ccAPI:          ccAPI,
		logger:         logger.Session("cc-service-broker-client", lager.Data{"cc-api": ccAPI}),
	}
}

//Create creates a service and returns an error if it fails
func (sb *ServiceBroker) Create(name, url, username, password string) error {
	log := sb.logger.Session("create-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

	path := "/v2/service_brokers"
	body := &BrokerValues{Name: name, BrokerURL: url, AuthUsername: username, AuthPassword: password}

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

//Update updates a service
func (sb *ServiceBroker) Update(serviceBrokerGUID, name, url, username, password string) error {
	log := sb.logger.Session("update-broker", lager.Data{"name": name, "url": url})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	path := fmt.Sprintf("/v2/service_brokers/%s", serviceBrokerGUID)
	body := BrokerValues{Name: name, BrokerURL: url, AuthUsername: username, AuthPassword: password}

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

//UpdateAll updates all services with the given url and username to the given password
func (sb *ServiceBroker) UpdateAll(url, username, password string) error {
	log := sb.logger.Session("update-all-brokers", lager.Data{"url": url})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

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
		if resource.Value.BrokerURL == url && resource.Value.AuthUsername == username {
			err = sb.Update(resource.Meta.GUID, resource.Value.Name, url, username, password)
			if err != nil {
				return err
			}
		} else {
			log.Debug("skipping-broker", lager.Data{"broker": resource.Value})
		}
	}

	log.Info("finished-cc-request")
	return nil
}

//EnableServiceAccess enables service access for the service having serviceName
func (sb *ServiceBroker) EnableServiceAccess(serviceName string) error {
	log := sb.logger.Session("enableservice-access", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccAPI, log)

	err := sp.Update(serviceName)
	if err != nil {
		return err
	}

	log.Debug("finished")

	return nil
}

//GetServiceBrokerGUIDByName obtains the broker guid corresponding to the passed name
func (sb *ServiceBroker) GetServiceBrokerGUIDByName(name string) (string, error) {
	log := sb.logger.Session("get-service-broker-guid-by-name", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/v2/service_brokers?q=name:%s", name)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	log.Debug("preparing-request", lager.Data{"path": path, "headers": headers})

	headers["Authorization"] = token

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

	guid := resources.Resources[0].Meta.GUID
	log.Debug("found", lager.Data{"service-broker-guid": guid})

	return guid, nil
}

//CheckServiceNameExists checks if a service with the passed name is already defined
func (sb *ServiceBroker) CheckServiceNameExists(name string) (bool, error) {
	exist := false
	log := sb.logger.Session("check-service-name-exists", lager.Data{"name": name})
	log.Debug("starting")

	token, err := sb.tokenGenerator.GetToken()
	if err != nil {
		log.Error("check-service-name-exists", err)
		return false, err
	}

	sp := NewServicePlan(sb.client, sb.tokenGenerator, sb.ccAPI, log)
	guid, err := sp.GetServiceGUIDByLabel(name, token)
	if err != nil {
		log.Error("get-service-guid-by-label", err)
		return false, err
	}
	if guid != "" {
		exist = true
	}
	log.Debug(fmt.Sprintf("check service name %s exists complete - returning %t", name, exist))

	return exist, nil
}

//CheckServiceInstancesExist checks if a service instance with the passed name is already registered
func (sb *ServiceBroker) CheckServiceInstancesExist(serviceName string) bool {
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

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	for _, plan := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s/service_instances", plan.Values.GUID)

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
func (sb *ServiceBroker) Delete(name string) error {
	log := sb.logger.Session("delete-broker", lager.Data{"name": name})
	log.Debug("starting")

	guid, err := sb.GetServiceBrokerGUIDByName(name)
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
