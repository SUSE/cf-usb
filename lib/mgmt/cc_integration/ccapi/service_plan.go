package ccapi

import (
	"encoding/json"
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
	Public bool `json:"public"`
}

type PlanResources struct {
	Resources []PlanResource `json:"resources"`
}

type PlanResource struct {
	Values PlanMetadata `json:"metadata"`
}

type PlanMetadata struct {
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

func (sp *ServicePlan) Update(serviceGuid string) error {
	token, err := sp.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	servicePlanGuids, err := sp.getServicePlanGuids(serviceGuid, token)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	body := PlanValues{Public: true}
	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sp.logger.Debug("update-service-plans", lager.Data{"service guid": serviceGuid})

	for _, servicePlanGuid := range servicePlanGuids {
		path := fmt.Sprintf("/v2/service_plans/%s", servicePlanGuid)
		request := httpclient.Request{Verb: "PUT", Endpoint: sp.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

		_, err = sp.client.Request(request)
		if err != nil {
			return err
		}

		sp.logger.Debug("update-service-plan", lager.Data{"service plan guid": servicePlanGuid})
	}

	return nil
}

func (sp *ServicePlan) getServicePlanGuids(serviceGuid, token string) ([]string, error) {
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

	var guids []string
	for _, value := range resources.Resources {
		guids = append(guids, value.Values.Guid)
	}

	sp.logger.Debug("get-service-plan", lager.Data{"service plan guids": guids})

	return guids, nil
}
