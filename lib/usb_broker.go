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
	configProvider config.ConfigProvider
	logger         lager.Logger
}

func NewUsbBroker(configProvider config.ConfigProvider, logger lager.Logger) *UsbBroker {
	return &UsbBroker{configProvider: configProvider, logger: logger.Session("usb-broker")}
}

func (broker *UsbBroker) Services() brokerapi.CatalogResponse {
	var catalog []brokerapi.Service
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		broker.logger.Fatal("retrive-configuration", err)
	}

	broker.logger.Info("get-catalog", lager.Data{})

	for _, driver := range config.Drivers {
		for _, instance := range driver.DriverInstances {
			service := instance.Service
			for _, dial := range instance.Dials {
				service.Plans = append(service.Plans, dial.Plan)
			}
			catalog = append(catalog, service)

		}
	}
	return brokerapi.CatalogResponse{Services: catalog}
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails, acceptsIncomplete bool) (brokerapi.ProvisioningResponse, bool, error) {
	broker.logger.Info("provision", lager.Data{"instanceID": instanceID})

	driver, err := broker.getDriver(serviceDetails.ServiceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}

	if driver == nil {
		return brokerapi.ProvisioningResponse{}, false, errors.New(fmt.Sprintf("Cannot find driver for %s", serviceDetails.ServiceID))
	}

	instance, err := driver.GetInstance(instanceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}
	if instance.Status == status.Exists {
		return brokerapi.ProvisioningResponse{}, false, brokerapi.ErrInstanceAlreadyExists
	}

	instance, err = driver.ProvisionInstance(instanceID, serviceDetails.PlanID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}
	broker.logger.Info("provision", lager.Data{"provisioned": instance.InstanceID})

	if instance.Status == status.Created {
		return brokerapi.ProvisioningResponse{}, false, nil
	}
	if instance.Status == status.InProgress {
		if !acceptsIncomplete {
			// TODO: clean up instance
			// driver.DeprovisionInstance(instanceID)
			return brokerapi.ProvisioningResponse{}, false, brokerapi.ErrAsyncRequired
		}

		return brokerapi.ProvisioningResponse{}, true, nil
	}

	return brokerapi.ProvisioningResponse{}, false, errors.New("Unknown instance state")
}

func (broker *UsbBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
	return false, brokerapi.ErrInstanceNotUpdateable
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
	driver, err := broker.getDriver(deprovisionDetails.ServiceID)
	if err != nil {
		return false, err
	}
	if driver == nil {
		return false, errors.New(fmt.Sprintf("Cannot find driver for %s", deprovisionDetails.ServiceID))
	}
	instance, err := driver.GetInstance(instanceID)
	if err != nil {
		return false, err
	}
	if instance.Status == status.DoesNotExist {
		return false, brokerapi.ErrInstanceDoesNotExist
	}

	instance, err = driver.DeprovisionInstance(instanceID)
	if err != nil {
		return false, err
	}
	broker.logger.Info("deprovision", lager.Data{"driver-response": instance.InstanceID})

	if instance.Status == status.Deleted {
		return false, nil
	}

	//  TODO:
	//  Is InProgress applicabale to deprovision?
	//  Should there be another state, e.g. status.DeprovisionInProgress vs status.ProvisionInProgrss
	if instance.Status == status.InProgress {
		return true, nil
	}

	return false, errors.New("Unknown instance state")
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
	var response brokerapi.BindingResponse

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

	response.Credentials, err = driver.GenerateCredentials(instanceID, bindingID)
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

func (broker *UsbBroker) LastOperation(instanceID string) (brokerapi.LastOperationResponse, error) {
	// TODO: how to get the driver for a instanceID. NOTE: the broker API does not require
	// the client to inclide the serviceID in the request
	driver, driverFound, err := broker.getDriverForServiceInstanceId(instanceID)
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}
	if !driverFound {
		return brokerapi.LastOperationResponse{}, brokerapi.ErrInstanceDoesNotExist
	}

	instance, err := driver.GetInstance(instanceID)
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}
	if instance.Status == status.DoesNotExist {
		return brokerapi.LastOperationResponse{}, brokerapi.ErrInstanceDoesNotExist
	}

	if instance.Status == status.Created || instance.Status == status.Exists {
		return brokerapi.LastOperationResponse{State: brokerapi.LastOperationSucceeded}, nil
	}

	if instance.Status == status.InProgress {
		return brokerapi.LastOperationResponse{State: brokerapi.LastOperationInProgress}, nil
	}

	if instance.Status == status.Error {
		return brokerapi.LastOperationResponse{State: brokerapi.LastOperationFailed}, nil
	}

	// TODO: what about instance.Status == status.Deleted ?
	return brokerapi.LastOperationResponse{}, errors.New("Unknown instance state")
}

func (broker *UsbBroker) getDriver(serviceID string) (*DriverProvider, error) {
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, driver := range config.Drivers {
		for driverInstanceID, driverInstance := range driver.DriverInstances {
			if driverInstance.Service.ID == serviceID {
				driverProvider := NewDriverProvider(driver.DriverType,
					broker.configProvider, driverInstanceID, broker.logger)
				return driverProvider, nil
			}
		}
	}

	return nil, errors.New("Driver not found")
}

func (broker *UsbBroker) getDriverForServiceInstanceId(instanceID string) (*DriverProvider, bool, error) {
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return nil, false, err
	}

	for _, driver := range config.Drivers {
		for driverInstanceID, _ := range driver.DriverInstances {
			driverProvider := NewDriverProvider(driver.DriverType,
				broker.configProvider, driverInstanceID, broker.logger)

			instance, err := driverProvider.GetInstance(instanceID)
			if err != nil {
				return nil, false, err
			}
			if instance.Status == status.DoesNotExist {
				continue
			}

			return driverProvider, true, nil
		}
	}

	return nil, false, nil
}
