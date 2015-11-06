package dummyasync

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/dummy/driverdata"
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
	return dummyAsyncDriver{logger: logger}
}

func (d dummyAsyncDriver) init(config *json.RawMessage) (DummyAsyncServiceConfig, error) {
	d.logger.Info("init-driver")

	d.logger.Info("init-driver", lager.Data{"configValue": string(*config)})
	dsp := DummyAsyncServiceConfig{}
	err := json.Unmarshal(*config, &dsp)
	if err != nil {
		return dsp, err
	}

	d.logger.Info("init-driver", lager.Data{"succed_count": dsp.SucceedCout})

	return dsp, err

}

func (d dummyAsyncDriver) Ping(request *json.RawMessage, response *bool) error {
	_, err := d.init(request)

	if err != nil {
		return err
	}

	*response = true
	return nil
}

func (d dummyAsyncDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d dummyAsyncDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("driver-get-config-schema")
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d dummyAsyncDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID})
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
	d.logger.Info("get-instance-request", lager.Data{"instanceID": request})

	response.Status = status.Created
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	serviceFilePath := filepath.Join(wd, request.InstanceID)
	if _, err := os.Stat(serviceFilePath); os.IsNotExist(err) {
		response.Status = status.DoesNotExist
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
	d.logger.Info("generate-credentials-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})

	*response = DummyAsyncServiceBindResponse{
		Content: "content",
	}

	return nil
}

func (d dummyAsyncDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.DoesNotExist
	d.logger.Info("credentials-exists-request", lager.Data{"instanceID": request.InstanceID,
		"credentialsID": request.CredentialsID})
	response.Status = status.DoesNotExist
	return nil
}

func (d dummyAsyncDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})

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
