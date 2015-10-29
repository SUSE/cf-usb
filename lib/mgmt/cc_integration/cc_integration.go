package cc_integration

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/cc_api"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaa_api"
	"github.com/pivotal-golang/lager"
)

type CCIntegrationInterface interface {
	Init() error
	CreateServiceBroker(name, url, username, password string) error
	UpdateServiceBroker(serviceBrokerGuid, name, url, username, password string) error
	EnableServiceAccess(servicePlanGuid, name, description, serviceGuid string, free, public bool) error
}

type CCIntegration struct {
	config         *config.Config
	tokenGenerator uaa_api.GetTokenInterface
	client         httpclient.HttpClient
	ccConfig       *CloudController
	logger         lager.Logger
}

type CloudController struct {
	Api               string `json:"api"`
	SkipTslValidation bool   `json:"skip_tsl_validation"`
}

func NewCCIntegration(apiConfig *config.Config, logger lager.Logger) CCIntegrationInterface {
	return &CCIntegration{
		config:         apiConfig,
		tokenGenerator: nil,
		client:         nil,
		ccConfig:       nil,
		logger:         logger,
	}
}

func (cci *CCIntegration) Init() error {
	conf := (*json.RawMessage)(cci.config.ManagementAPI.CloudController)
	cc := CloudController{}
	err := json.Unmarshal(*conf, &cc)
	if err != nil {
		return err
	}
	cci.ccConfig = &cc
	cci.client = httpclient.NewHttpClient(cc.SkipTslValidation)

	info := cc_api.NewGetInfo(cci.ccConfig.Api, cci.client, cci.logger)
	tokenUrl, err := info.GetTokenEndpoint()
	if err != nil {
		return err
	}

	tokenGenerator := uaa_api.NewTokenGenerator(tokenUrl, cci.config.ManagementAPI.UaaClient, cci.config.ManagementAPI.UaaSecret, cci.client)
	cci.tokenGenerator = tokenGenerator

	cci.logger.Info("init-uaa-token-generator", lager.Data{"token url": tokenUrl, "skip tls validation": cc.SkipTslValidation})

	return nil
}

func (cci *CCIntegration) CreateServiceBroker(name, url, username, password string) error {
	sb := cc_api.NewServiceBroker(cci.client, cci.tokenGenerator, cci.ccConfig.Api, cci.logger)

	err := sb.Create(name, url, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (cci *CCIntegration) UpdateServiceBroker(serviceBrokerGuid, name, url, username, password string) error {
	sb := cc_api.NewServiceBroker(cci.client, cci.tokenGenerator, cci.ccConfig.Api, cci.logger)

	err := sb.Update(serviceBrokerGuid, name, url, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (cci *CCIntegration) EnableServiceAccess(servicePlanGuid, name, description, serviceGuid string, free, public bool) error {
	sp := cc_api.NewServicePlan(cci.client, cci.tokenGenerator, cci.ccConfig.Api, cci.logger)

	err := sp.Update(servicePlanGuid, name, description, serviceGuid, free, public)
	if err != nil {
		return err
	}

	return nil
}
