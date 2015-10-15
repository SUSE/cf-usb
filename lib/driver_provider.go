package lib

import (
	"errors"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
	"github.com/xeipuuv/gojsonschema"
)

type DriverProvider struct {
	driverType string
	client     *rpc.Client

	DriverProperties config.DriverProperties
}

func NewDriverProvider(driverType string, driverProperties config.DriverProperties, logger lager.Logger) (*DriverProvider, error) {
	provider := DriverProvider{}

	driverPath := os.Getenv("USB_DRIVER_PATH")

	if driverPath == "" {
		driverPath = "drivers"
	}
	driverPath = filepath.Join(driverPath, driverType)
	if runtime.GOOS == "windows" {
		driverPath = driverPath + ".exe"
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {

		return &provider, err
	}
	provider.client = client
	provider.DriverProperties = driverProperties
	provider.driverType = driverType

	configSchema, err := provider.GetConfigSchema()
	if err != nil {
		return &provider, err
	}

	schemaLoader := gojsonschema.NewStringLoader(configSchema)
	configLoader := gojsonschema.NewGoLoader(driverProperties.DriverConfiguration)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return &provider, err
	}
	if !result.Valid() {
		err = errors.New("Invalid configuration schema")

		errData := lager.Data{}
		for _, e := range result.Errors() {
			errData[e.Field()] = e.Description()
		}
		logger.Error("driver-init", err, errData)
		return &provider, err
	}

	initRequest := model.DriverInitRequest{}
	initRequest.DriverConfig = driverProperties.DriverConfiguration
	//TODO: Implement dails
	_, err = provider.Init(initRequest)

	return &provider, err
}

func (p *DriverProvider) Init(request model.DriverInitRequest) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Init", p.driverType), request, &result)
	return result, err
}

func (p *DriverProvider) Ping() (bool, error) {
	result := false
	err := p.client.Call(fmt.Sprintf("%s.Ping", p.driverType), "", &result)

	return result, err
}

func (p *DriverProvider) GetDailsSchema() (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.GetDailsSchema", p.driverType), "", &result)
	return result, err
}

func (p *DriverProvider) GetConfigSchema() (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.GetConfigSchema", p.driverType), "", &result)
	return result, err
}

func (p *DriverProvider) ProvisionInstance(provisonRequest model.ProvisionInstanceRequest) (bool, error) {
	var result bool
	err := p.client.Call(fmt.Sprintf("%s.ProvisionInstance", p.driverType), provisonRequest, &result)
	return result, err
}

func (p *DriverProvider) InstanceExists(instanceID string) (bool, error) {
	var result bool
	err := p.client.Call(fmt.Sprintf("%s.InstanceExists", p.driverType), instanceID, &result)
	return result, err
}

func (p *DriverProvider) GenerateCredentials(credentialsRequest model.CredentialsRequest) (interface{}, error) {
	var result interface{}
	err := p.client.Call(fmt.Sprintf("%s.GenerateCredentials", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) CredentialsExist(credentialsRequest model.CredentialsRequest) (bool, error) {
	var result bool
	err := p.client.Call(fmt.Sprintf("%s.CredentialsExist", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) RevokeCredentials(credentialsRequest model.CredentialsRequest) (interface{}, error) {
	var result interface{}
	err := p.client.Call(fmt.Sprintf("%s.RevokeCredentials", p.driverType), credentialsRequest, &result)
	return result, err
}

func (p *DriverProvider) DeprovisionInstance(deprovisionRequest string) (interface{}, error) {
	var result interface{}
	err := p.client.Call(fmt.Sprintf("%s.DeprovisionInstance", p.driverType), deprovisionRequest, &result)
	return result, err
}
