package driver

import (
	"fmt"
	"log"
)

type dummyDriver struct {
	Driver
}

func NewDummyDriver() Driver {
	return dummyDriver{}
}

func (driver dummyDriver) Provision(request interface{}, response *interface{}) error {
	log.Println("i am provisioning!!!!")
	*response = fmt.Sprintf("I am provisioning with %s", request)
	return nil

}
func (driver dummyDriver) Deprovision(request string, response *string) error {
	return nil
}
func (driver dummyDriver) Bind(request string, response *string) error {
	return nil
}
func (driver dummyDriver) Unbind(request string, response *string) error {
	return nil
}
func (driver dummyDriver) Update(request string, response *string) error {
	return nil
}
func (driver dummyDriver) GetCatalog(request string, response *string) error {
	*response = "Dummy catalog"
	return nil
}
func (driver dummyDriver) GetInstances(request string, response *string) error {
	return nil
}
