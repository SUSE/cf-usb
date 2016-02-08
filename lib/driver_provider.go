package lib

import (
	"fmt"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/pivotal-golang/lager"
)

type DriverProvider struct {
	driverType string
	driverPath string

	logger           lager.Logger
	ConfigProvider   config.ConfigProvider
	driverInstanceID string
}

func NewDriverProvider(driversPath string, driverType string, configProvider config.ConfigProvider,
	driverInstanceID string, logger lager.Logger) *DriverProvider {
	p := DriverProvider{}

	p.ConfigProvider = configProvider
	p.driverInstanceID = driverInstanceID
	p.driverType = driverType
	p.logger = logger.Session("driver-provider")

	p.driverPath = getDriverPath(driversPath, driverType)
	p.logger.Debug("new-driver-provider", lager.Data{"service-method": fmt.Sprintf("%s.DriverType", p.driverType), "driver-path": p.driverPath})
	return &p
}

func (p *DriverProvider) ProvisionInstance(instanceID, planID string) (driver.Instance, error) {
	p.logger.Debug("provision-instance-request", lager.Data{"instance-id": instanceID, "plan-id": planID})

	var result driver.Instance

	driverInstance, err := p.ConfigProvider.LoadDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	provisonRequest := driver.ProvisionInstanceRequest{}
	provisonRequest.Config = driverInstance.Configuration
	provisonRequest.InstanceID = instanceID

	for _, d := range driverInstance.Dials {
		if d.Plan.ID == planID {
			provisonRequest.Dials = d.Configuration
			break
		}
	}

	p.logger.Debug("provision-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.ProvisionInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.ProvisionInstance", p.driverType), p.driverPath,
		provisonRequest, &result)

	return result, err
}

func (p *DriverProvider) GetInstance(instanceID string) (driver.Instance, error) {
	p.logger.Debug("get-instance-request", lager.Data{"instance-id": instanceID})

	var result driver.Instance

	driverInstance, _, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	instanceRequest := driver.GetInstanceRequest{}
	instanceRequest.Config = driverInstance.Configuration
	instanceRequest.InstanceID = instanceID

	p.logger.Debug("get-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GetInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GetInstance", p.driverType),
		p.driverPath, instanceRequest, &result)

	return result, err
}

func (p *DriverProvider) GenerateCredentials(instanceID, credentialsID string) (interface{}, error) {
	p.logger.Debug("generate-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result interface{}

	driverInstance, _, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.GenerateCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("generate-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GenerateCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GenerateCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) GetCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	p.logger.Debug("get-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result driver.Credentials

	driverInstance, _, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.GetCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("get-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GetCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GetCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) RevokeCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	p.logger.Debug("revoke-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result driver.Credentials

	driverInstance, _, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.RevokeCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("revoke-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.RevokeCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.RevokeCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) DeprovisionInstance(instanceID string) (driver.Instance, error) {
	p.logger.Debug("deprovision-instance-request", lager.Data{"instance-id": instanceID})

	var result driver.Instance

	driverInstance, _, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	deprovisionRequest := driver.DeprovisionInstanceRequest{}
	deprovisionRequest.Config = driverInstance.Configuration
	deprovisionRequest.InstanceID = instanceID

	p.logger.Debug("deprovision-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.DeprovisionInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.DeprovisionInstance", p.driverType),
		p.driverPath, deprovisionRequest, &result)

	return result, err
}
