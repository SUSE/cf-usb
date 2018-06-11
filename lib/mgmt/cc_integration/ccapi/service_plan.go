package ccapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager"
)

//ServicePlanInterface defines a service plan actions
type ServicePlanInterface interface {
	Update(ServiceGUID) error
	GetServiceGUIDByLabel(ServiceName, uaaapi.BearerToken) (ServiceGUID, error)
	GetServicePlans(ServiceGUID, uaaapi.BearerToken) (*PlanResources, error)
}

//ServicePlan holds details for ServicePlan
type ServicePlan struct {
	client         httpclient.HTTPClient
	tokenGenerator uaaapi.GetTokenInterface
	ccAPI          string
	logger         lager.Logger
}

//PlanResources holds the resources for the plan
type PlanResources struct {
	Resources []struct {
		Metadata struct {
			GUID PlanGUID `json:"guid"`
		} `json:"metadata"`
		Entity struct {
			Name        string `json:"name"`
			Free        bool   `json:"free"`
			Description string `json:"description"`
			Public      bool   `json:"public"`
			ServiceGUID string `json:"service_guid"`
		} `json:"entity"`
	} `json:"resources"`
}

// A PlanGUID is the unique identifier for a service plan
type PlanGUID string

//ServiceResources holds the resources for the service
type ServiceResources struct {
	Resources []struct {
		Metadata struct {
			GUID ServiceGUID `json:"guid"`
		} `json:"metadata"`
	} `json:"resources"`
}

// A ServiceName is the name of a service (~= sidecar deployment)
type ServiceName string

// A ServiceGUID is the GUID of a service (~= sidecar deployment)
type ServiceGUID string

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
func (sp *ServicePlan) Update(serviceGUID ServiceGUID) error {
	log := sp.logger.Session("update-service-plans", lager.Data{"service-broker": serviceGUID})
	log.Debug("starting")
	defer log.Debug("finished")

	token, err := sp.tokenGenerator.GetToken()
	if err != nil {
		return err
	}

	servicePlans, err := sp.GetServicePlans(serviceGUID, token)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Authorization": string(token),
	}

	log.Debug("initializing")

	for _, value := range servicePlans.Resources {
		path := fmt.Sprintf("/v2/service_plans/%s", value.Metadata.GUID)

		body := value.Entity
		body.Public = true
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

//GetServiceGUIDByLabel returns a cloud controller service GUID by its label
func (sp *ServicePlan) GetServiceGUIDByLabel(serviceLabel ServiceName, token uaaapi.BearerToken) (ServiceGUID, error) {
	log := sp.logger.Session("get-service-guid-by-label", lager.Data{"service-label": serviceLabel})
	log.Debug("starting")
	defer log.Debug("finished")

	path := fmt.Sprintf("/v2/services?q=label:%s", serviceLabel)

	headers := map[string]string{
		"Authorization": string(token),
		"Content-Type":  "application/x-www-form-urlencoded; charset=UTF-8",
		"Accept":        "application/json; charset=utf-8",
	}

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
		return "", nil
	}

	guid := resources.Resources[0].Metadata.GUID
	log.Debug("found", lager.Data{"service-guid": guid})

	return guid, nil
}

//GetServicePlans returns the service plans for a cloud controller service GUID and token
func (sp *ServicePlan) GetServicePlans(serviceGUID ServiceGUID, token uaaapi.BearerToken) (*PlanResources, error) {
	log := sp.logger.Session("get-service-plans", lager.Data{"service-guid": serviceGUID})
	log.Debug("starting")
	defer log.Debug("finished")

	path := fmt.Sprintf("/v2/service_plans?q=service_guid:%s", serviceGUID)

	headers := map[string]string{
		"Authorization": string(token),
		"Content-Type":  "application/x-www-form-urlencoded; charset=UTF-8",
		"Accept":        "application/json; charset=utf-8",
	}

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
