package passthrough

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/passthrough/config"
	"github.com/hpcloud/cf-usb/driver/passthrough/driverdata"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type passthroughDriver struct {
	logger            lager.Logger
	stateFileLocation string
}

func NewPassthroughDriver(logger lager.Logger) driver.Driver {
	return passthroughDriver{
		logger:            logger.Session("passthrough-driver"),
		stateFileLocation: "/s/data/passthrough_driver.json",
	}
}

func (d passthroughDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

	var c config.CredentialsConfig

	*response = false

	err := json.Unmarshal(*request, &c)
	if err != nil {
		return err
	}

	var resp map[string]string

	err = json.Unmarshal([]byte(c.StaticConfig), &resp)
	if err != nil {
		return err
	}

	*response = true

	return nil
}

func (d passthroughDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d passthroughDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d passthroughDriver) GetParametersSchema(request string, response *string) error {
	d.logger.Info("get-parameters-schema-request", lager.Data{"request": request})

	//Does not support custom parameters

	return nil
}

func (d passthroughDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	state, err := d.getState()
	if err != nil {
		response.Status = status.Error
		return err
	}

	instance := config.ServiceInstance{
		ID:          request.InstanceID,
		Credentials: make([]*config.ServiceCredential, 0),
	}

	state.Instances = append(state.Instances, &instance)

	err = d.saveState(state)
	if err != nil {
		response.Status = status.Error
		return err
	}

	response.Status = status.Created

	return nil
}

func (d passthroughDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

	state, err := d.getState()
	if err != nil {
		response.Status = status.Error
		return err
	}

	for _, instance := range state.Instances {
		if instance.ID == request.InstanceID {
			response.InstanceID = instance.ID
			response.Status = status.Exists
			return nil
		}
	}

	response.Status = status.DoesNotExist
	return nil
}

func (d passthroughDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	state, err := d.getState()
	if err != nil {
		return err
	}

	cred := config.ServiceCredential{
		ID: request.CredentialsID,
	}

	for _, instance := range state.Instances {
		if instance.ID == request.InstanceID {
			instance.Credentials = append(instance.Credentials, &cred)
			break
		}
	}

	err = d.saveState(state)
	if err != nil {
		return err
	}

	var c config.CredentialsConfig

	err = json.Unmarshal(*request.Config, &c)
	if err != nil {
		return err
	}

	var resp map[string]string

	err = json.Unmarshal([]byte(c.StaticConfig), &resp)
	if err != nil {
		return err
	}

	*response = resp

	return nil
}

func (d passthroughDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	state, err := d.getState()
	if err != nil {
		return err
	}

	for _, instance := range state.Instances {
		if instance.ID == request.InstanceID {
			for _, credential := range instance.Credentials {
				if credential.ID == request.CredentialsID {
					response.CredentialsID = credential.ID
					response.Status = status.Exists
					return nil
				}
			}
		}
	}

	response.Status = status.DoesNotExist
	return nil
}

func (d passthroughDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	state, err := d.getState()
	if err != nil {
		response.Status = status.Error
		return err
	}

	for _, instance := range state.Instances {
		if instance.ID == request.InstanceID {
			credIndex := -1
			for i, credential := range instance.Credentials {
				if credential.ID == request.CredentialsID {
					credIndex = i
					break
				}
			}
			if credIndex > -1 {
				instance.Credentials = append(instance.Credentials[:credIndex], instance.Credentials[credIndex+1:]...)

				err = d.saveState(state)
				if err != nil {
					response.Status = status.Error
					return err
				}

				response.Status = status.Deleted
				return nil
			}
		}
	}

	response.Status = status.DoesNotExist
	return nil
}

func (d passthroughDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	state, err := d.getState()
	if err != nil {
		response.Status = status.Error
		return err
	}

	instanceIndex := -1
	for i, instance := range state.Instances {
		if instance.ID == request.InstanceID {
			instanceIndex = i
			break
		}
	}

	if instanceIndex > -1 {
		state.Instances = append(state.Instances[:instanceIndex], state.Instances[instanceIndex+1:]...)

		err = d.saveState(state)
		if err != nil {
			response.Status = status.Error
			return err
		}

		response.Status = status.Deleted
		return nil
	}

	response.Status = status.DoesNotExist

	return nil
}

func (d passthroughDriver) getState() (*config.DriverState, error) {

	var state config.DriverState
	stat, err := os.Stat(d.stateFileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			state = config.DriverState{
				Instances: make([]*config.ServiceInstance, 0),
			}
			return &state, nil
		}
	}

	if stat.Size() == 0 {
		state = config.DriverState{
			Instances: make([]*config.ServiceInstance, 0),
		}
		return &state, nil
	}

	file, err := ioutil.ReadFile(d.stateFileLocation)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (d passthroughDriver) saveState(state *config.DriverState) error {

	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(d.stateFileLocation, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
