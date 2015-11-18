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
		logger:         logger,
	}
}

func (sp *ServicePlan) Update(serviceName string) error {
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

	sp.logger.Debug("update-service-plans", lager.Data{"service guid": serviceGuid})

	for _, value := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s", value.Values.Guid)

		body := PlanValues{Name: value.Entity.Name, Free: value.Entity.Free, Description: value.Entity.Description, Public: value.Entity.Public, ServiceGuid: value.Entity.ServiceGuid}
		values, err := json.Marshal(body)
		if err != nil {
			return err
		}

		request := httpclient.Request{Verb: "PUT", Endpoint: sp.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

		_, err = sp.client.Request(request)
		if err != nil {
			return err
		}

		sp.logger.Debug("update-service-plan", lager.Data{"service plan guid": value.Values.Guid})
	}

	return nil
}

func (sp *ServicePlan) getServiceGuidByLabel(serviceLabel, token string) (string, error) {
	path := fmt.Sprintf("/v2/services?q=label:%s", serviceLabel)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

	response, err := sp.client.Request(findRequest)
	if err != nil {
		return "", err
	}

	resources := &ServiceResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return "", err
	}
	sp.logger.Debug("result", lager.Data{"response": string(response)})

	if len(resources.Resources) == 0 {
		return "", errors.New(fmt.Sprintf("Service %s not found", serviceLabel))
	}

	guid := resources.Resources[0].Values.Guid

	sp.logger.Debug("get-service", lager.Data{"service guid by name": serviceLabel})

	return guid, nil
}

func (sp *ServicePlan) getServicePlans(serviceGuid, token string) (*PlanResources, error) {
	path := fmt.Sprintf("/v2/service_plans?q=service_guid:%s", serviceGuid)

	headers := make(map[string]string)
	headers["Authorization"] = token
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["Accept"] = "application/json; charset=utf-8"

	findRequest := httpclient.Request{Verb: "GET", Endpoint: sp.ccApi, ApiUrl: path, Headers: headers, StatusCode: 200}

	response, err := sp.client.Request(findRequest)
	if err != nil {
		return nil, err
	}

	resources := &PlanResources{}
	err = json.Unmarshal(response, resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
