package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hpcloud/cf-usb/lib/config"
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

func (driver *dummyDriver) Provision(request interface{}, response *interface{}) error {
	log.Println("i am provisioning!!!!")
	*response = fmt.Sprintf("I am provisioning with %s", request)
	return nil

}
func (driver *dummyDriver) Deprovision(request string, response *string) error {
	return nil
}
func (driver *dummyDriver) Bind(request string, response *string) error {
	return nil
}
func (driver *dummyDriver) Unbind(request string, response *string) error {
	return nil
}
func (driver *dummyDriver) Update(request string, response *string) error {
	return nil
}
