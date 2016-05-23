package ccapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager"
)

//ServicePlanInterface defines a service plan actions
type ServicePlanInterface interface {
	Update(serviceGUID string) error
	GetServiceGUIDByLabel(string, string) (string, error)
	GetServicePlans(string, string) (*PlanResources, error)
}

//ServicePlan holds details for ServicePlan
type ServicePlan struct {
	client         httpclient.HTTPClient
	tokenGenerator uaaapi.GetTokenInterface
	ccAPI          string
	logger         lager.Logger
}

//PlanValues describes the plan values - json mapped
type PlanValues struct {
	Name        string `json:"name"`
	Free        bool   `json:"free"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	ServiceGUID string `json:"service_guid"`
}

//PlanResources holds the resources for the plan
type PlanResources struct {
	Resources []PlanResource `json:"resources"`
}

//PlanResource defines a resource to be used in a plan, mapped to json
type PlanResource struct {
	Values PlanMetadata `json:"metadata"`
	Entity PlanValues   `json:"entity"`
}

//PlanMetadata holds the GUID metadata - mapped to json
type PlanMetadata struct {
	GUID string `json:"guid"`
}

//ServiceResources holds the resources for the service
type ServiceResources struct {
	Resources []ServiceResource `json:"resources"`
}

//ServiceResource defines  resource to be used by a service
type ServiceResource struct {
	Values BrokerMetadata `json:"metadata"`
}

//NewServicePlan instantiates and returns a service plan
func NewServicePlan(client httpclient.HTTPClient, token uaaapi.GetTokenInterface, ccAPI string, logger lager.Logger) ServicePlanInterface {
	return &ServicePlan{
		client:         client,
		tokenGenerator: token,
		ccAPI:          ccAPI,
		logger:         logger.Session("cc-service-plans-client", lager.Data{"cc-api": ccAPI}),
	}
}

//Update updates a service to use this plan
func (sp *ServicePlan) Update(serviceName string) error {
	log := sp.logger.Session("update-service-plans", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	token, err := sp.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	serviceGUID, err := sp.GetServiceGUIDByLabel(serviceName, token)
	if err != nil {
		return err
	}

	servicePlans, err := sp.GetServicePlans(serviceGUID, token)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	log.Debug("initializing", lager.Data{"service guid": serviceGUID})

	for _, value := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s", value.Values.GUID)

		body := PlanValues{Name: value.Entity.Name, Free: value.Entity.Free, Description: value.Entity.Description, Public: true, ServiceGUID: value.Entity.ServiceGUID}
		values, err := json.Marshal(body)
		if err != nil {
			return err
		}

		log.Debug("preparing-request", lager.Data{"request-content": string(values)})

		request := httpclient.Request{Verb: "PUT", Endpoint: sp.ccAPI, APIURL: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

		log.Info("starting-cc-request", lager.Data{"path": path})

		_, err = sp.client.Request(request)
		if err != nil {
			return err
		}

		log.Info("finished-cc-request")
	}

	return nil
}

//GetServiceGUIDByLabel returns a service GUID from its label
func (sp *ServicePlan) GetServiceGUIDByLabel(serviceLabel, token string) (string, error) {
	log := sp.logger.Session("get-service-guid-by-label", lager.Data{"service-label": serviceLabel})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/services?q=label:%s", serviceLabel)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccAPI, APIURL: path, Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path, "verb": "GET"})

	response, err := sp.client.Request(findRequest)
	if err != nil {
		return "", err
	}

	log.Debug("cc-reponse", lager.Data{"response": string(response)})
	log.Info("finished-cc-request")

	resources := &ServiceResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return "", err
	}

	if len(resources.Resources) == 0 {
		return "", fmt.Errorf("Service %s not found", serviceLabel)
	}

	guid := resources.Resources[0].Values.GUID
	log.Debug("found", lager.Data{"service-guid": guid})

	return guid, nil
}

//GetServicePlans returns the service plans for a serviceGUID and token
func (sp *ServicePlan) GetServicePlans(serviceGUID, token string) (*PlanResources, error) {
	log := sp.logger.Session("get-service-plans", lager.Data{"service-guid": serviceGUID})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/service_plans?q=service_guid:%s", serviceGUID)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccAPI, APIURL: path, Headers: headers, StatusCode: 200}

	log.Info("starting-cc-request", lager.Data{"path": path, "verb": "GET"})

	response, err := sp.client.Request(findRequest)
	if err != nil {
		return nil, err
	}

	log.Debug("cc-reponse", lager.Data{"response": string(response)})
	log.Info("finished-cc-request")

	resources := &PlanResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
