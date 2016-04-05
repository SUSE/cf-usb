package dummy

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
	return dummyDriver{logger: logger.Session("dummy-driver")}
}

func (d dummyDriver) init(config *json.RawMessage) (DummyServiceConfig, error) {
	d.logger.Info("init-driver", lager.Data{"configValue": string(*config)})

	dsp := DummyServiceConfig{}
	err := json.Unmarshal(*config, &dsp)
	if err != nil {
		return DummyServiceConfig{}, err
	}

	return dsp, nil
}

func (d dummyDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

	_, err := d.init(request)

	if err != nil {
		return err
	}

	*response = true

	return nil
}

func (d dummyDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d dummyDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")

	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d dummyDriver) GetParametersSchema(request string, response *string) error {
	d.logger.Info("get-parameters-schema-request", lager.Data{"request": request})

	parametersSchema, err := driverdata.Asset("schemas/parameters.json")

	if err != nil {
		return err
	}

	*response = string(parametersSchema)

	return nil
}

func (d dummyDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	response.Status = status.Created
	response.InstanceID = request.InstanceID
	response.Description = "Instance created"

	return nil
}

func (d dummyDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

	response.Status = status.DoesNotExist
	if request.InstanceID == "instanceID" {
		response.Status = status.Exists
	}

	return nil
}

func (d dummyDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	*response = DummyServiceBindResponse{
		UserName: "user",
		Password: "pass",
	}

	return nil
}

func (d dummyDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	response.Status = status.DoesNotExist

	if request.CredentialsID == "credentialsID" {
		response.Status = status.Exists
	}

	return nil
}

func (d dummyDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	response.Status = status.Deleted

	return nil
}

func (d dummyDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	response.Status = status.Deleted
	response.InstanceID = request.InstanceID
	response.Description = "Instance deleted"

	return nil
}
