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

type ServicePlanInterface interface {
	Update(serviceGuid string) error
}

type ServicePlan struct {
	client         httpclient.HttpClient
	tokenGenerator uaaapi.GetTokenInterface
	ccApi          string
	logger         lager.Logger
}

type PlanValues struct {
	Name        string `json:"name"`
	Free        bool   `json:"free"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	ServiceGuid string `json:"service_guid"`
}

type PlanResources struct {
	Resources []PlanResource `json:"resources"`
}

type PlanResource struct {
	Values PlanMetadata `json:"metadata"`
	Entity PlanValues   `json:"entity"`
}

type PlanMetadata struct {
	Guid string `json:"guid"`
}

type ServiceResources struct {
	Resources []ServiceResource `json:"resources"`
}

type ServiceResource struct {
	Values BrokerMetadata `json:"metadata"`
}

type ServiceMetadata struct {
	Guid string `json:"guid"`
}

func NewServicePlan(client httpclient.HttpClient, token uaaapi.GetTokenInterface, ccApi string, logger lager.Logger) ServicePlanInterface {
	return &ServicePlan{
		client:         client,
		tokenGenerator: token,
		ccApi:          ccApi,
		logger:         logger.Session("cc-service-plans-client", lager.Data{"cc-api": ccApi}),
	}
}

func (sp *ServicePlan) Update(serviceName string) error {
	log := sp.logger.Session("update-service-plans", lager.Data{"service-name": serviceName})
	log.Debug("starting")

	token, err := sp.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	serviceGuid, err := sp.getServiceGuidByLabel(serviceName, token)
	if err != nil {
		return err
	}

	servicePlans, err := sp.getServicePlans(serviceGuid, token)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	log.Debug("initializing", lager.Data{"service guid": serviceGuid})

	for _, value := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s", value.Values.Guid)

		body := PlanValues{Name: value.Entity.Name, Free: value.Entity.Free, Description: value.Entity.Description, Public: true, ServiceGuid: value.Entity.ServiceGuid}
		values, err := json.Marshal(body)
		if err != nil {
			return err
		}

		log.Debug("preparing-request", lager.Data{"request-content": string(values)})

		request := httpclient.Request{Verb: "PUT", Endpoint: sp.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

		log.Info("starting-cc-request", lager.Data{"path": path})

		_, err = sp.client.Request(request)
		if err != nil {
			return err
		}

		log.Info("finished-cc-request")
	}

	return nil
}

func (sp *ServicePlan) getServiceGuidByLabel(serviceLabel, token string) (string, error) {
	log := sp.logger.Session("get-service-guid-by-label", lager.Data{"service-label": serviceLabel})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/services?q=label:%s", serviceLabel)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

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
		return "", errors.New(fmt.Sprintf("Service %s not found", serviceLabel))
	}

	guid := resources.Resources[0].Values.Guid
	log.Debug("found", lager.Data{"service-guid": guid})

	return guid, nil
}

func (sp *ServicePlan) getServicePlans(serviceGuid, token string) (*PlanResources, error) {
	log := sp.logger.Session("get-service-plans", lager.Data{"service-guid": serviceGuid})
	log.Debug("starting")

	path := fmt.Sprintf("/v2/service_plans?q=service_guid:%s", serviceGuid)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

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
