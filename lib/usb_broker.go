package lib

import (
	"errors"
	"fmt"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/hpcloud/cf-usb/lib/config"
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
	for _, driverProvider := range broker.driverProverders {
		if driverProvider.Instance != nil {
			service := driverProvider.Instance.Service
			for _, dial := range driverProvider.Instance.Dials {
				service.Plans = append(service.Plans, dial.Plan)
			}
			catalog = append(catalog, service)
		}
	}
	return catalog
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails) error {
	broker.logger.Info("provision", lager.Data{"instanceID": instanceID})

	driver := broker.getDriver(serviceDetails.ID)
	if driver == nil {
		return errors.New(fmt.Sprintf("Cannot find driver for %s", serviceDetails.ID))
	}

	instance, err := driver.GetInstance(instanceID)
	if err != nil {
		return err
	}
	if instance.Status == status.Exists {
		return brokerapi.ErrInstanceAlreadyExists
	}

	//TODO: add async
	instance, err = driver.ProvisionInstance(instanceID, serviceDetails.PlanID)
	if err != nil {
		return err
	}
	broker.logger.Info("provision", lager.Data{"provisioned": instance.InstanceID})

	return nil
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails) error {
	driver := broker.getDriver(deprovisionDetails.ServiceID)
	if driver == nil {
		return errors.New(fmt.Sprintf("Cannot find driver for %s", deprovisionDetails.ServiceID))
	}
	instance, err := driver.GetInstance(instanceID)
	if err != nil {
		return err
	}
	if instance.Status == status.DoesNotExist {
		return brokerapi.ErrInstanceDoesNotExist
	}

	instance, err = driver.DeprovisionInstance(instanceID)
	if err != nil {
		return err
	}
	broker.logger.Info("deprovision", lager.Data{"driver-response": instance.InstanceID})
	return nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	var response interface{}

	driver := broker.getDriver(details.ServiceID)

	credentials, err := driver.GetCredentials(instanceID, bindingID)
	if err != nil {
		return response, err
	}
	if credentials.Status == status.Exists {
		return response, brokerapi.ErrBindingAlreadyExists
	}

	response, err = driver.GenerateCredentials(instanceID, bindingID)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	driver := broker.getDriver(details.ServiceID)

	credentials, err := driver.GetCredentials(instanceID, bindingID)
	if err != nil {
		return err
	}
	if credentials.Status == status.DoesNotExist {
		return brokerapi.ErrBindingDoesNotExist
	}

	credentials, err = driver.RevokeCredentials(instanceID, bindingID)
	if err != nil {
		return err
	}
	broker.logger.Info("unbind", lager.Data{"driver-response": credentials.Status})

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
