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
	configProvider   config.ConfigProvider
	logger           lager.Logger
}

func NewUsbBroker(drivers []*DriverProvider, configProvider config.ConfigProvider, logger lager.Logger) *UsbBroker {
	return &UsbBroker{driverProverders: drivers, configProvider: configProvider, logger: logger}
}

func (broker *UsbBroker) Services() []brokerapi.Service {
	var catalog []brokerapi.Service

	for _, driverProvider := range broker.driverProverders {
		driverInstance, err := broker.configProvider.LoadDriverInstance(driverProvider.driverInstanceID)
		if err != nil {
			broker.logger.Error("get-services", err, lager.Data{"driverInstanceID": driverProvider.driverInstanceID})
			continue
		}
		service := driverInstance.Service
		for _, dial := range driverInstance.Dials {
			service.Plans = append(service.Plans, dial.Plan)
		}
		catalog = append(catalog, service)

	}
	return catalog
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails) error {
	broker.logger.Info("provision", lager.Data{"instanceID": instanceID})

	driver, err := broker.getDriver(serviceDetails.ID)
	if err != nil {
		return err
	}

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
	driver, err := broker.getDriver(deprovisionDetails.ServiceID)
	if err != nil {
		return err
	}
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

	driver, err := broker.getDriver(details.ServiceID)
	if err != nil {
		return response, err
	}

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
	driver, err := broker.getDriver(details.ServiceID)
	if err != nil {
		return err
	}

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

func (broker *UsbBroker) getDriver(serviceID string) (*DriverProvider, error) {
	for _, driverProvider := range broker.driverProverders {
		service, err := broker.configProvider.GetService(driverProvider.driverInstanceID)
		if err != nil {
			return nil, err
		}
		if service.ID == serviceID {
			return driverProvider, nil
		}
	}
	return nil, nil
}
