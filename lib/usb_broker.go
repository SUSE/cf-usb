package lib

import (
	"log"

	"github.com/hpcloud/gocfbroker"
)

var usbBroker UsbBroker

type UsbBroker struct {
	driverProverders []DriverProvider
}

func NewUsbBroker(drivers []DriverProvider) *UsbBroker {
	return &UsbBroker{driverProverders: drivers}
}

func (broker *UsbBroker) Provision(instanceID string, req gocfbroker.ProvisionRequest) (gocfbroker.ProvisionResponse, error) {
	//TODO: Errors are not treated in the gocfbroker
	res := gocfbroker.ProvisionResponse{}
	driver := broker.getDriver(req.ServiceID)
	response, err := driver.Provision("req")
	if err != nil {
		return res, err
	}
	log.Println("response is", response)
	return res, nil
}

func (broker *UsbBroker) Deprovision(instanceID, serviceID, planID string) error {
	return nil
}

func (broker *UsbBroker) Update(instanceID string, req gocfbroker.UpdateProvisionRequest) error {
	return nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, req gocfbroker.BindingRequest) (gocfbroker.BindingResponse, error) {
	res := gocfbroker.BindingResponse{}
	return res, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID, serviceID, planID string) error {
	return nil

}

func (broker *UsbBroker) getDriver(serviceID string) DriverProvider {
	dp := DriverProvider{}
	for _, driverProvider := range broker.driverProverders {
		for _, s := range driverProvider.DriverProperties.Services {
			if serviceID == s.ID {
				return driverProvider
			}
		}
	}
	return dp
}
