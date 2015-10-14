package lib

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/natefinch/pie"
)

type DriverProvider struct {
	driverType string
	client     *rpc.Client

	DriverProperties config.DriverProperties
}

func NewDriverProvider(driverType string, driverProperties config.DriverProperties) (*DriverProvider, error) {
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

	_, err = provider.Init(driverProperties)

	return &provider, err
}

func (p *DriverProvider) Init(driverProperties config.DriverProperties) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Init", p.driverType), driverProperties, &result)
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

func (p *DriverProvider) ProvisionInstance(provisonRequest model.ProvisionInstanceRequest) (string, error) {
	var result string
	err := p.client.Call(fmt.Sprintf("%s.Provision", p.driverType), provisonRequest, &result)
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
