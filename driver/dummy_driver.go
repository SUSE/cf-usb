package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/hpcloud/gocfbroker"
)

type DummyServiceProperties struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type dummyDriver struct {
	driverProperties config.DriverProperties
	Driver
}

func NewDummyDriver() Driver {
	return &dummyDriver{}
}
func (driver *dummyDriver) Init(driverProperties config.DriverProperties, response *string) error {
	driver.driverProperties = driverProperties

	log.Println("Driver dummy initialized")
	for _, service := range driverProperties.Services {
		log.Println("Using serviceID:", service.ID)
		log.Println("Service Description", service.Description)
		for _, plan := range service.Plans {
			log.Println("PlanID:", plan.ID)
			log.Println("PlanName:", plan.Name)
		}
	}

	conf := (*json.RawMessage)(driverProperties.DriverConfiguration)
	log.Println(string(*conf))
	dsp := DummyServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	log.Println("property_one:", dsp.PropOne)
	log.Println("property_two:", dsp.PropTwo)

	*response = "Sucessfully initialized driver"
	return nil

}

func (driver *dummyDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	log.Println("Provisioning ", request.InstanceID)
	log.Println("with plan id", request.BrokerProvisionRequest.PlanID)
	*response = fmt.Sprintf("http://example-dashboard.com/9189kdfsk0vfnku")
	return nil

}
func (driver *dummyDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	log.Println("Deprovisioning ", request.InstanceID)
	log.Println("with plan id ", request.PlanID)
	*response = "Successfully deprovisoned"
	return nil
}

func (driver *dummyDriver) Update(request model.DriverUpdateRequest, response *string) error {
	log.Println("Updating", request.InstanceID)
	log.Println("with plan id ", request.BrokerUpdateRequest.PlanID)
	return nil
}

func (driver *dummyDriver) Bind(request model.DriverBindRequest, response *gocfbroker.BindingResponse) error {
	log.Println("Binding", request.InstanceID)
	log.Println("using planID", request.BrokerBindRequest.PlanID)
	log.Println("on appUD", request.BrokerBindRequest.AppGUID)

	data := []byte(`{"user": "testuser","password":"testpassword"}`)
	response.Credentials = (*json.RawMessage)(&data)
	response.SyslogDrainURL = "don't think this is used"
	return nil
}
func (driver *dummyDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	log.Println("Unbinding ", request.BindingID)
	log.Println("with planID", request.PlanID)
	return nil
}
