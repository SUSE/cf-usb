package lib

import "github.com/hpcloud/gocfbroker"

var usbBroker UsbBroker

type UsbBroker struct {
}

func (broker *UsbBroker) Provision(instanceID string, req gocfbroker.ProvisionRequest) (gocfbroker.ProvisionResponse, error) {
	res := gocfbroker.ProvisionResponse{}
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
