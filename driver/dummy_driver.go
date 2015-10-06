package driver

import (
	"encoding/json"
	"fmt"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type DummyServiceProperties struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type dummyDriver struct {
	driverProperties config.DriverProperties
	logger           lager.Logger
	Driver
}

func NewDummyDriver(logger lager.Logger) Driver {
	return &dummyDriver{logger: logger}
}
func (driver *dummyDriver) Init(driverProperties config.DriverProperties, response *string) error {
	driver.driverProperties = driverProperties

	driver.logger.Info("Driver dummy initialized")
	for _, service := range driverProperties.Services {
		driver.logger.Info("provision-service:", lager.Data{"serviceID": service.ID, "description": service.Description})
		for _, plan := range service.Plans {
			driver.logger.Info("pans", lager.Data{"PlanID": plan.ID, "PlanName": plan.Name})
		}
	}

	conf := (*json.RawMessage)(driverProperties.DriverConfiguration)
	driver.logger.Info(string(*conf))
	dsp := DummyServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	driver.logger.Info("serviceProperties", lager.Data{"property_one": dsp.PropOne, "property_two": dsp.PropTwo})

	*response = "Sucessfully initialized driver"
	return nil

}

func (driver *dummyDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	driver.logger.Info("Provisioning", lager.Data{"instance-id": request.InstanceID, "plan-id": request.ServiceDetails.PlanID})
	*response = fmt.Sprintf("http://example-dashboard.com/9189kdfsk0vfnku")

	if request.InstanceID == "exists" {
		return brokerapi.ErrInstanceAlreadyExists
	}

	return nil

}
func (driver *dummyDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	driver.logger.Info("deprovision-request", lager.Data{"instance-id": request.InstanceID})
	*response = "Successfully deprovisoned"
	return nil
}

func (driver *dummyDriver) Bind(request model.DriverBindRequest, response *json.RawMessage) error {
	driver.logger.Info("bind-request", lager.Data{"instanceID": request.InstanceID,
		"planID": request.BindDetails.PlanID, "appID": request.BindDetails.AppGUID})

	data := []byte(`{"user": "testuser","password":"testpassword"}`)
	response = (*json.RawMessage)(&data)
	return nil
}
func (driver *dummyDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	driver.logger.Info("unbind-request", lager.Data{"bindingID": request.BindingID, "InstanceID": request.InstanceID})
	return nil
}
