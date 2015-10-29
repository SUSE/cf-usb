package lib

import (
	"errors"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
	"github.com/xeipuuv/gojsonschema"
)

type DriverProvider struct {
	driverType string
	driverPath string

	logger   lager.Logger
	Instance *config.DriverInstance
}

func NewDriverProvider(driverType string, instance *config.DriverInstance, logger lager.Logger) *DriverProvider {
	p := DriverProvider{}

	p.Instance = instance
	p.driverType = driverType
	p.logger = logger

	driverPath := os.Getenv("USB_DRIVER_PATH")
	if driverPath == "" {
		driverPath = "drivers"
	}

	driverPath = filepath.Join(driverPath, driverType)
	if runtime.GOOS == "windows" {
		driverPath = driverPath + ".exe"
	}
	p.driverPath = driverPath

	return &p
}

func (p *DriverProvider) Validate() error {
	client, err := p.createProviderClient()
	if err != nil {
		return err
	}
	defer client.Close()

	err = p.validateConfigSchema(client)
	if err != nil {
		return err
	}

	err = p.validateDialsSchema(client)
	if err != nil {
		return err
	}

	pong, err := p.Ping()
	if err != nil {
		return err
	}

	if !pong {
		err = errors.New("Cannot reach server.")
		return err
	}

	return nil
}

func (p *DriverProvider) Ping() (bool, error) {
	result := false
	err := p.createClientAndCall(fmt.Sprintf("%s.Ping", p.driverType), "", &result)
	return result, err
}

func (p *DriverProvider) GetDailsSchema() (string, error) {
	return p.createClientAndInvoke(p.getDailsSchema)
}

func (p *DriverProvider) GetConfigSchema() (string, error) {
	return p.createClientAndInvoke(p.getConfigSchema)
}

func (p *DriverProvider) ProvisionInstance(instanceID, planID string) (driver.Instance, error) {
	var result driver.Instance
	provisonRequest := driver.ProvisionInstanceRequest{}
	provisonRequest.Config = p.Instance.Configuration
	provisonRequest.InstanceID = instanceID
	for _, d := range p.Instance.Dials {
		if d.Plan.ID == planID {
			provisonRequest.Dials = d.Configuration
			break
		}
	}

	err := p.createClientAndCall(fmt.Sprintf("%s.ProvisionInstance", p.driverType), provisonRequest, &result)
	return result, err
}

func (p *DriverProvider) GetInstance(instanceID string) (driver.Instance, error) {
	var result driver.Instance
	instanceRequest := driver.GetInstanceRequest{}
	instanceRequest.Config = p.Instance.Configuration
	instanceRequest.InstanceID = instanceID

	err := p.createClientAndCall(fmt.Sprintf("%s.GetInstance", p.driverType), instanceRequest, &result)
	return result, err
}

func (p *DriverProvider) GenerateCredentials(instanceID, credentialsID string) (interface{}, error) {
	var result interface{}
	credentialsRequest := driver.GenerateCredentialsRequest{}
	credentialsRequest.Config = p.Instance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	err := p.createClientAndCall(fmt.Sprintf("%s.GenerateCredentials", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) GetCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	var result driver.Credentials
	credentialsRequest := driver.GetCredentialsRequest{}
	credentialsRequest.Config = p.Instance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	err := p.createClientAndCall(fmt.Sprintf("%s.GetCredentials", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) RevokeCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	var result driver.Credentials
	credentialsRequest := driver.RevokeCredentialsRequest{}
	credentialsRequest.Config = p.Instance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID
	err := p.createClientAndCall(fmt.Sprintf("%s.RevokeCredentials", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) DeprovisionInstance(instanceID string) (driver.Instance, error) {
	var result driver.Instance
	deprovisionRequest := driver.DeprovisionInstanceRequest{}
	deprovisionRequest.Config = p.Instance.Configuration
	deprovisionRequest.InstanceID = instanceID
	err := p.createClientAndCall(fmt.Sprintf("%s.DeprovisionInstance", p.driverType), deprovisionRequest, &result)
	return result, err
}

func (p *DriverProvider) createProviderClient() (*rpc.Client, error) {
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, p.driverPath)
	return client, err
}

func (p *DriverProvider) validateDialsSchema(client *rpc.Client) error {
	dialSchema, err := p.getDailsSchema(client)

	if err != nil {
		return err
	}

	dialsSchemaLoader := gojsonschema.NewStringLoader(dialSchema)
	for _, dial := range p.Instance.Dials {
		dialLoader := gojsonschema.NewGoLoader(dial.Configuration)
		result, err := gojsonschema.Validate(dialsSchemaLoader, dialLoader)
		if err != nil {
			return err
		}
		if !result.Valid() {
			err = errors.New("Invalid dials configuration")

			errData := lager.Data{}
			for _, e := range result.Errors() {
				errData[e.Field()] = e.Description()
			}
			p.logger.Error("driver-init", err, errData)
			return err
		}
	}

	return nil
}

func (p *DriverProvider) validateConfigSchema(client *rpc.Client) error {
	configSchema, err := p.getConfigSchema(client)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewStringLoader(configSchema)
	configLoader := gojsonschema.NewGoLoader(p.Instance.Configuration)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		err = errors.New("Invalid configuration schema")

		errData := lager.Data{}
		for _, e := range result.Errors() {
			errData[e.Field()] = e.Description()
		}
		p.logger.Error("driver-init", err, errData)
		return err
	}

	return nil
}

func (p *DriverProvider) createClientAndCall(serviceMethod string, args interface{}, reply interface{}) error {
	client, err := p.createProviderClient()
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Call(serviceMethod, args, reply)
}

func (p *DriverProvider) createClientAndInvoke(call func(*rpc.Client) (string, error)) (string, error) {
	client, err := p.createProviderClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	return call(client)
}

func (p *DriverProvider) getDailsSchema(client *rpc.Client) (string, error) {
	var result string
	err := client.Call(fmt.Sprintf("%s.GetDailsSchema", p.driverType), "", &result)
	return result, err
}

func (p *DriverProvider) getConfigSchema(client *rpc.Client) (string, error) {
	var result string
	err := client.Call(fmt.Sprintf("%s.GetConfigSchema", p.driverType), "", &result)
	return result, err
}
