package cc_api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaa_api"
	"github.com/pivotal-golang/lager"
)

type ServicePlanInterface interface {
	Update(servicePlanGuid, name, description, serviceGuid string, free, public bool) error
}

type ServicePlan struct {
	client         httpclient.HttpClient
	tokenGenerator uaa_api.GetTokenInterface
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

func NewServicePlan(client httpclient.HttpClient, token uaa_api.GetTokenInterface, ccApi string, logger lager.Logger) ServicePlanInterface {
	return &ServicePlan{
		client:         client,
		tokenGenerator: token,
		ccApi:          ccApi,
		logger:         logger,
	}
}

func (sp *ServicePlan) Update(servicePlanGuid, name, description, serviceGuid string, free, public bool) error {
	path := fmt.Sprintf("/v2/service_plans/%s", servicePlanGuid)
	body := PlanValues{Name: name, Free: free, Description: description, Public: public, ServiceGuid: serviceGuid}

	values, err := json.Marshal(body)
	if err != nil {
		return err
	}

	sp.logger.Debug("update-service-plan", lager.Data{"service plan name": name})

	token, err := sp.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Authorization"] = token

	request := httpclient.Request{Verb: "PUT", Endpoint: sp.ccApi, ApiUrl: path, Body: strings.NewReader(string(values)), Headers: headers, StatusCode: 201}

	_, err = sp.client.Request(request)
	if err != nil {
		return err
	}

	return nil
	//return sp.createUpdateRequest("PUT", sp.ccApi, path, token, strings.NewReader(string(values)), 201)
}

//func (sp *ServicePlan) createUpdateRequest(verb, endpoint, apiUrl, accessToken string, body io.ReadSeeker, statusCode int) error {
//	request, err := http.NewRequest(verb, endpoint+apiUrl, body)
//	if err != nil {
//		return errors.New("Error building request")
//	}

//	if accessToken != "" {
//		request.Header.Set("Authorization", accessToken)
//	}

//	response, err := sp.client.Do(request)
//	if err != nil {
//		return err
//	}

//	defer response.Body.Close()

//	responseBody, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return err
//	}

//	if response.StatusCode != statusCode {
//		return errors.New(fmt.Sprintf("status code: %d, body: %s", response.StatusCode, responseBody))
//	}

//	return nil
//}
