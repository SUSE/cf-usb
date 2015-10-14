package dummydriver

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type DummyServiceProperties struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type DummyServiceBindResponse struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type dummyDriver struct {
	driverProperties config.DriverProperties
	logger           lager.Logger
}

func NewDummyDriver(logger lager.Logger) driver.Driver {
	return &dummyDriver{logger: logger}
}
func (driver *dummyDriver) Init(driverProperties config.DriverProperties, response *string) error {
	driver.logger.Info("init-driver")
	driver.driverProperties = driverProperties

	for _, service := range driverProperties.Services {
		driver.logger.Info("init-driver", lager.Data{"serviceID": service.ID, "description": service.Description})
		for _, plan := range service.Plans {
			driver.logger.Info("pans", lager.Data{"PlanID": plan.ID, "PlanName": plan.Name})
		}
	}

	conf := (*json.RawMessage)(driverProperties.DriverConfiguration)
	driver.logger.Info("init-driver", lager.Data{"configValue": string(*conf)})
	dsp := DummyServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	driver.logger.Info("init-driver", lager.Data{"property_one": dsp.PropOne, "property_two": dsp.PropTwo})

	*response = "Sucessfully initialized driver"
	return nil

}

func (driver *dummyDriver) Ping(request string, response *bool) error {
	if request == "exists" {
		*response = true
	}
	*response = false
	return nil
}

func (driver *dummyDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := data.Asset("schemas/dails.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (driver *dummyDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := data.Asset("scehmas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (driver *dummyDriver) ProvisionInstance(request model.ProvisionInstanceRequest, response *bool) error {
	driver.logger.Info("Provisioning", lager.Data{"instance-id": request.InstanceID})
	*response = true

	if request.InstanceID == "instanceID" {
		return brokerapi.ErrInstanceAlreadyExists
	}

	return nil
}

func (driver *dummyDriver) InstanceExists(request string, response *bool) error {
	if request == "instanceID" {
		*response = true
	}
	*response = false

	return nil
}

func (driver *dummyDriver) GenerateCredentials(request model.CredentialsRequest, response *interface{}) error {
	driver.logger.Info("generate-credentials-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})

	*response = DummyServiceBindResponse{
		UserName: "user",
		Password: "pass",
	}

	return nil
}

func (driver *dummyDriver) CredentialsExist(request model.CredentialsRequest, response *bool) error {
	driver.logger.Info("credentials-exists-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})
	if request.CredentialsID == "credentialsID" {
		*response = true
	}

	*response = false
	return nil
}

func (driver *dummyDriver) RevokeCredentials(request model.CredentialsRequest, response *interface{}) error {
	driver.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})

	*response = ""
	return nil
}

func (driver *dummyDriver) DeprovisionInstance(request string, response *interface{}) error {
	driver.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	*response = "Successfully deprovisoned"
	return nil
}
