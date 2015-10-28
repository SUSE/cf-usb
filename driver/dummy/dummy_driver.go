package dummydriver

import (
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/dummy/driverdata"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type DummyServiceConfig struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type DummyServiceBindResponse struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type dummyDriver struct {
	logger lager.Logger
}

func NewDummyDriver(logger lager.Logger) driver.Driver {
	return dummyDriver{logger: logger}
}
func (d dummyDriver) init(config *json.RawMessage) (DummyServiceConfig, error) {
	d.logger.Info("init-driver")

	d.logger.Info("init-driver", lager.Data{"configValue": string(*config)})
	dsp := DummyServiceConfig{}
	err := json.Unmarshal(*config, &dsp)
	if err != nil {
		return dsp, err
	}

	d.logger.Info("init-driver", lager.Data{"property_one": dsp.PropOne, "property_two": dsp.PropTwo})

	return dsp, err

}

func (d dummyDriver) Ping(request *json.RawMessage, response *bool) error {
	_, err := d.init(request)

	if err != nil {
		return err
	}

	*response = true
	return nil
}

func (d dummyDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d dummyDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("driver-get-config-schema")
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d dummyDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID})
	response.Status = status.Created

	return nil
}

func (d dummyDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instanceID": request})
	response.Status = status.DoesNotExist
	if request.InstanceID == "instanceID" {
		response.Status = status.Exists
	}

	return nil
}

func (d dummyDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})

	*response = DummyServiceBindResponse{
		UserName: "user",
		Password: "pass",
	}

	return nil
}

func (d dummyDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.DoesNotExist
	d.logger.Info("credentials-exists-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})
	if request.CredentialsID == "credentialsID" {
		response.Status = status.Exists
	}

	return nil
}

func (d dummyDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})

	response.Status = status.Deleted
	return nil
}

func (d dummyDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	response.Status = status.Deleted

	return nil
}
