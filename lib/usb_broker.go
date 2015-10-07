package lib

import (
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

	driverProvisionRequest := model.DriverProvisionRequest{
		InstanceID:     instanceID,
		ServiceDetails: serviceDetails,
	}

	driverResponse, err := driver.Provision(driverProvisionRequest)
	if err != nil {
		return err
	}
	broker.logger.Info("provision", lager.Data{"driver-response": driverResponse})

	return nil
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails) error {
	driver := broker.getDriver(deprovisionDetails.ServiceID)

	driverDeprovisionRequest := model.DriverDeprovisionRequest{
		InstanceID: instanceID,
	}

	response, err := driver.Deprovision(driverDeprovisionRequest)
	if err != nil {
		return err
	}
	broker.logger.Info("deprovision", lager.Data{"driver-response": response})
	return nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (interface{}, error) {
	var response interface{}

	driver := broker.getDriver(details.ServiceID)

	driverBindRequest := model.DriverBindRequest{
		InstanceID:  instanceID,
		BindingID:   bindingID,
		BindDetails: details,
	}

	response, err := driver.Bind(driverBindRequest)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	driver := broker.getDriver(details.ServiceID)

	driverUnbindRequest := model.DriverUnbindRequest{
		InstanceID: instanceID,
		BindingID:  bindingID,
	}

	response, err := driver.Unbind(driverUnbindRequest)
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
