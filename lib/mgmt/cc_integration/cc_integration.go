package ccintegration

import (
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager"
)

type CCIntegrationInterface interface {
	CreateServiceBroker(name string) error
	UpdateServiceBroker(name string) error
	EnableServicesAccess() error
}

type CCIntegration struct {
	config         *config.Config
	tokenGenerator uaaapi.GetTokenInterface
	client         httpclient.HttpClient
	logger         lager.Logger
}

func NewCCIntegration(config *config.Config, tokenGenerator uaaapi.GetTokenInterface, client httpclient.HttpClient, logger lager.Logger) CCIntegrationInterface {
	return &CCIntegration{
		config:         config,
		tokenGenerator: tokenGenerator,
		client:         client,
		logger:         logger,
	}
}

func (cci *CCIntegration) CreateServiceBroker(name string) error {
	sb := ccapi.NewServiceBroker(cci.client, cci.tokenGenerator, cci.config.ManagementAPI.CloudController.Api, cci.logger)

	err := sb.Create(name, cci.config.BrokerAPI.ExternalUrl, cci.config.BrokerAPI.Credentials.Username, cci.config.BrokerAPI.Credentials.Password)
	if err != nil {
		return err
	}

	return nil
}

func (cci *CCIntegration) UpdateServiceBroker(name string) error {
	sb := ccapi.NewServiceBroker(cci.client, cci.tokenGenerator, cci.config.ManagementAPI.CloudController.Api, cci.logger)

	err := sb.Update(name, cci.config.BrokerAPI.ExternalUrl, cci.config.BrokerAPI.Credentials.Username, cci.config.BrokerAPI.Credentials.Password)
	if err != nil {
		return err
	}

	return nil
}

func (cci *CCIntegration) EnableServicesAccess() error {
	sp := ccapi.NewServicePlan(cci.client, cci.tokenGenerator, cci.config.ManagementAPI.CloudController.Api, cci.logger)

	for _, d := range cci.config.Drivers {
		for _, di := range d.DriverInstances {

			err := sp.Update(di.Service.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
