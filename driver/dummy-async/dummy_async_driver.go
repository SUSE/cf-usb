package dummyasync

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/dummy-async/driverdata"
	"github.com/hpcloud/cf-usb/driver/status"

	"github.com/pivotal-golang/lager"
)

type DummyAsyncServiceConfig struct {
	SucceedCout string `json:"succeed_count"`
}

type DummyAsyncServiceBindResponse struct {
	Content string `json:"content"`
}

type dummyAsyncDriver struct {
	logger lager.Logger
}

func NewDummyAsyncDriver(logger lager.Logger) driver.Driver {
	return dummyAsyncDriver{logger: logger.Session("async-dummy-driver")}
}

func (d dummyAsyncDriver) init(config *json.RawMessage) (DummyAsyncServiceConfig, error) {
	d.logger.Info("init-driver", lager.Data{"configValue": string(*config)})

	dsp := DummyAsyncServiceConfig{}
	err := json.Unmarshal(*config, &dsp)
	if err != nil {
		return dsp, err
	}

	d.logger.Info("init-driver", lager.Data{"succeeded-count": dsp.SucceedCout})

	return dsp, err

}

func (d dummyAsyncDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

	_, err := d.init(request)

	if err != nil {
		return err
	}

	*response = true
	return nil
}

func (d dummyAsyncDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d dummyAsyncDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d dummyAsyncDriver) GetParametersSchema(request string, response *string) error {
	//Does not support custom parameters
	return nil
}

func (d dummyAsyncDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	config, err := d.init(request.Config)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	serviceFilePath := filepath.Join(wd, request.InstanceID)

	ioutil.WriteFile(serviceFilePath, []byte(config.SucceedCout), 0644)
	response.Status = status.InProgress

	return nil
}

func (d dummyAsyncDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

	response.Status = status.Created
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	serviceFilePath := filepath.Join(wd, request.InstanceID)
	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		response.Status = status.DoesNotExist
		return nil
	}

	content, err := ioutil.ReadFile(serviceFilePath)
	if err != nil {
		return err
	}

	step, err := strconv.Atoi(string(content))
	if err != nil {
		return err
	}

	if step > 0 {
		value := strconv.Itoa(step - 1)
		ioutil.WriteFile(serviceFilePath, []byte(value), 0644)
		response.Status = status.InProgress
	}

	return nil
}

func (d dummyAsyncDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	*response = DummyAsyncServiceBindResponse{
		Content: "content",
	}

	return nil
}

func (d dummyAsyncDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	response.Status = status.DoesNotExist

	if request.CredentialsID == "credentialsID" {
		response.Status = status.Exists
	}

	return nil
}

func (d dummyAsyncDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	response.Status = status.Deleted

	return nil
}

func (d dummyAsyncDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	serviceFilePath := filepath.Join(wd, request.InstanceID)

	os.Remove(serviceFilePath)
	response.Status = status.Deleted

	return nil
}
