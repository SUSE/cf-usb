package lib

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

var usbBroker UsbBroker

type UsbBroker struct {
	driverProverders []*DriverProvider
	brokerConfig     *config.Config
	logger           lager.Logger
}

func NewUsbBroker(drivers []*DriverProvider, config *config.Config, logger lager.Logger) *UsbBroker {
	return &UsbBroker{driverProverders: drivers, brokerConfig: config, logger: logger}
}

func (broker *UsbBroker) Services() []brokerapi.Service {
	return broker.brokerConfig.ServiceCatalog
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails) error {
	broker.logger.Info("provision", lager.Data{"instanceID": instanceID})

	driver := broker.getDriver(serviceDetails.ID)

	driverProvisionRequest := model.ProvisionInstanceRequest{
		InstanceID: instanceID,
		Dails:      json.RawMessage{},
	}

	driverResponse, err := driver.ProvisionInstance(driverProvisionRequest)
	if err != nil {
		return err
	}
	broker.logger.Info("provision", lager.Data{"driver-response": driverResponse})

	return nil
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails) error {
	driver := broker.getDriver(deprovisionDetails.ServiceID)

	response, err := driver.DeprovisionInstance(instanceID)
	if err != nil {
		return err
	}
	broker.logger.Info("deprovision", lager.Data{"driver-response": response})
	return nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	var response interface{}

	driver := broker.getDriver(details.ServiceID)

	driverCredentialsRequest := model.CredentialsRequest{
		InstanceID:    instanceID,
		CredentialsID: bindingID,
	}

	response, err := driver.GenerateCredentials(driverCredentialsRequest)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	driver := broker.getDriver(details.ServiceID)

	driverCredentialsRequest := model.CredentialsRequest{
		InstanceID:    instanceID,
		CredentialsID: bindingID,
	}

	response, err := driver.RevokeCredentials(driverCredentialsRequest)
	if err != nil {
		return err
	}
	broker.logger.Info("unbind", lager.Data{"driver-response": response})

	return nil

}

func (broker *UsbBroker) getDriver(serviceID string) *DriverProvider {

	for _, driverProvider := range broker.driverProverders {
		for _, s := range driverProvider.DriverProperties.Services {
			if serviceID == s.ID {
				return driverProvider
			}
		}
	}
	return nil
}
