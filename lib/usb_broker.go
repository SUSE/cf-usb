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
	var catalog []brokerapi.Service
	// TODO: advertise only services that have a driver running
	for _, driverProvider := range broker.driverProverders {
		service := driverProvider.Instance.Service
		for _, dial := range driverProvider.Instance.Dials {
			service.Plans = append(service.Plans, dial.Plan)
		}
		catalog = append(catalog, driverProvider.Instance.Service)
	}
	return catalog
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails) error {
	broker.logger.Info("provision", lager.Data{"instanceID": instanceID})

	driver := broker.getDriver(serviceDetails.ID)

	exists, err := driver.InstanceExists(instanceID)
	if err != nil {
		return err
	}
	if exists {
		return brokerapi.ErrInstanceAlreadyExists
	}

	var dialDetails json.RawMessage
	for _, dial := range driver.Instance.Dials {
		if dial.Plan.ID == serviceDetails.PlanID {
			dialDetails = *dial.Configuration
			break
		}
	}

	driverProvisionRequest := model.ProvisionInstanceRequest{
		InstanceID: instanceID,
		Dails:      dialDetails,
	}

	created, err := driver.ProvisionInstance(driverProvisionRequest)
	if err != nil {
		return err
	}
	broker.logger.Info("provision", lager.Data{"provisioned": created})

	return nil
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails) error {
	driver := broker.getDriver(deprovisionDetails.ServiceID)

	exists, err := driver.InstanceExists(instanceID)
	if err != nil {
		return err
	}
	if !exists {
		return brokerapi.ErrInstanceDoesNotExist
	}

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

	exists, err := driver.CredentialsExist(driverCredentialsRequest)
	if err != nil {
		return response, err
	}
	if exists {
		return response, brokerapi.ErrBindingAlreadyExists
	}

	response, err = driver.GenerateCredentials(driverCredentialsRequest)
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

	exists, err := driver.CredentialsExist(driverCredentialsRequest)
	if err != nil {
		return err
	}
	if !exists {
		return brokerapi.ErrBindingDoesNotExist
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
		if driverProvider.Instance.Service.ID == serviceID {
			return driverProvider
		}
	}
	return nil
}
