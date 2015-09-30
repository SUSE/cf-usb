package lib

import (
	"log"

	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/hpcloud/gocfbroker"
)

var usbBroker UsbBroker

type UsbBroker struct {
	driverProverders []*DriverProvider
}

func NewUsbBroker(drivers []*DriverProvider) *UsbBroker {
	return &UsbBroker{driverProverders: drivers}
}

func (broker *UsbBroker) Provision(instanceID string, request gocfbroker.ProvisionRequest) (gocfbroker.ProvisionResponse, error) {

	response := gocfbroker.ProvisionResponse{}

	driver := broker.getDriver(request.ServiceID)

	driverProvisionRequest := model.DriverProvisionRequest{
		InstanceID:             instanceID,
		BrokerProvisionRequest: request,
	}

	driverResponse, err := driver.Provision(driverProvisionRequest)
	if err != nil {
		return response, err
	}
	log.Println("Privision response:", driverResponse)
	response.DashboardURL = driverResponse

	return response, nil
}

func (broker *UsbBroker) Deprovision(instanceID, serviceID, planID string) error {

	driver := broker.getDriver(serviceID)

	driverDeprovisionRequest := model.DriverDeprovisionRequest{
		InstanceID: instanceID,
		PlanID:     planID,
	}

	response, err := driver.Deprovision(driverDeprovisionRequest)
	if err != nil {
		return err
	}
	log.Println("Deprovision response", response)
	return nil
}

func (broker *UsbBroker) Update(instanceID string, request gocfbroker.UpdateProvisionRequest) error {
	driver := broker.getDriver(request.ServiceID)

	driverUpdateRequest := model.DriverUpdateRequest{
		InstanceID:          instanceID,
		BrokerUpdateRequest: request,
	}

	response, err := driver.Update(driverUpdateRequest)
	if err != nil {
		return err
	}

	log.Println("Update response:", response)
	return nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, request gocfbroker.BindingRequest) (gocfbroker.BindingResponse, error) {
	response := gocfbroker.BindingResponse{}

	driver := broker.getDriver(request.ServiceID)

	driverBindRequest := model.DriverBindRequest{
		InstanceID:        instanceID,
		BindingID:         bindingID,
		BrokerBindRequest: request,
	}

	response, err := driver.Bind(driverBindRequest)
	if err != nil {
		return response, nil
	}

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID, serviceID, planID string) error {
	driver := broker.getDriver(serviceID)

	driverUnbindRequest := model.DriverUnbindRequest{
		InstanceID: instanceID,
		BindingID:  bindingID,
		PlanID:     planID,
	}

	response, err := driver.Unbind(driverUnbindRequest)
	if err != nil {
		return err
	}
	log.Println("Unbind response:", response)
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
