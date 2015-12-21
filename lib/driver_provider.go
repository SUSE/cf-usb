package lib

import (
	"encoding/json"
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

	logger           lager.Logger
	ConfigProvider   config.ConfigProvider
	driverInstanceID string
}

func NewDriverProvider(driverType string, configProvider config.ConfigProvider,
	driverInstanceID string, logger lager.Logger) *DriverProvider {
	p := DriverProvider{}

	p.ConfigProvider = configProvider
	p.driverInstanceID = driverInstanceID
	p.driverType = driverType
	p.logger = logger.Session("driver-provider")
	p.driverPath = getDriverPath(p.driverType)

	return &p
}

func (p *DriverProvider) ProvisionInstance(instanceID, planID string) (driver.Instance, error) {
	p.logger.Debug("provision-instance-request", lager.Data{"instance-id": instanceID, "plan-id": planID})

	var result driver.Instance

	driverInstance, err := p.ConfigProvider.LoadDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	provisonRequest := driver.ProvisionInstanceRequest{}
	provisonRequest.Config = driverInstance.Configuration
	provisonRequest.InstanceID = instanceID

	for _, d := range driverInstance.Dials {
		if d.Plan.ID == planID {
			provisonRequest.Dials = d.Configuration
			break
		}
	}

	p.logger.Debug("provision-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.ProvisionInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.ProvisionInstance", p.driverType), p.driverPath,
		provisonRequest, &result)

	return result, err
}

func (p *DriverProvider) GetInstance(instanceID string) (driver.Instance, error) {
	p.logger.Debug("get-instance-request", lager.Data{"instance-id": instanceID})

	var result driver.Instance

	driverInstance, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	instanceRequest := driver.GetInstanceRequest{}
	instanceRequest.Config = driverInstance.Configuration
	instanceRequest.InstanceID = instanceID

	p.logger.Debug("get-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GetInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GetInstance", p.driverType),
		p.driverPath, instanceRequest, &result)

	return result, err
}

func (p *DriverProvider) GenerateCredentials(instanceID, credentialsID string) (interface{}, error) {
	p.logger.Debug("generate-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result interface{}

	driverInstance, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.GenerateCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("generate-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GenerateCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GenerateCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) GetCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	p.logger.Debug("get-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result driver.Credentials

	driverInstance, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.GetCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("get-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.GetCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.GetCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) RevokeCredentials(instanceID, credentialsID string) (driver.Credentials, error) {
	p.logger.Debug("revoke-credentials-request", lager.Data{"instance-id": instanceID, "credentials-id": credentialsID})

	var result driver.Credentials

	driverInstance, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	credentialsRequest := driver.RevokeCredentialsRequest{}
	credentialsRequest.Config = driverInstance.Configuration
	credentialsRequest.CredentialsID = credentialsID
	credentialsRequest.InstanceID = instanceID

	p.logger.Debug("revoke-credentials-call-client", lager.Data{"service-method": fmt.Sprintf("%s.RevokeCredentials", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.RevokeCredentials", p.driverType),
		p.driverPath, credentialsRequest, &result)

	return result, err
}

func (p *DriverProvider) DeprovisionInstance(instanceID string) (driver.Instance, error) {
	p.logger.Debug("deprovision-instance-request", lager.Data{"instance-id": instanceID})

	var result driver.Instance

	driverInstance, err := p.ConfigProvider.GetDriverInstance(p.driverInstanceID)
	if err != nil {
		return result, err
	}

	deprovisionRequest := driver.DeprovisionInstanceRequest{}
	deprovisionRequest.Config = driverInstance.Configuration
	deprovisionRequest.InstanceID = instanceID

	p.logger.Debug("deprovision-instance-call-client", lager.Data{"service-method": fmt.Sprintf("%s.DeprovisionInstance", p.driverType), "driver-path": p.driverPath})

	err = createClientAndCall(fmt.Sprintf("%s.DeprovisionInstance", p.driverType),
		p.driverPath, deprovisionRequest, &result)

	return result, err
}

func Validate(driverInstance config.DriverInstance, driverType string, logger lager.Logger) error {
	log := logger.Session("validate-driver-instance", lager.Data{"name": driverInstance.Name, "type": driverType})

	client, err := createProviderClient(getDriverPath(driverType))
	if err != nil {
		return err
	}
	defer client.Close()

	log.Debug("validate-config-schema", lager.Data{"configuration": string(*driverInstance.Configuration)})

	err = validateConfigSchema(client, driverType, driverInstance.Configuration, logger)
	if err != nil {
		return err
	}

	log.Debug("validate-dials-schema", lager.Data{"dials-count": len(driverInstance.Dials)})

	err = validateDialsSchema(client, driverType, driverInstance, logger)
	if err != nil {
		return err
	}

	log.Debug("ping-driver", lager.Data{"configuration": string(*driverInstance.Configuration)})

	pong, err := ping(driverInstance.Configuration, driverType)
	if err != nil {
		return err
	}

	if !pong {
		err = errors.New("Cannot reach server.")
		return err
	}

	return nil
}

func ping(configuration *json.RawMessage, driverType string) (bool, error) {
	result := false
	driverPath := getDriverPath(driverType)
	err := createClientAndCall(fmt.Sprintf("%s.Ping", driverType), driverPath, configuration, &result)

	return result, err
}

func createProviderClient(driverPath string) (*rpc.Client, error) {
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)

	return client, err
}

func validateDialsSchema(client *rpc.Client, driverType string, driverInstance config.DriverInstance,
	logger lager.Logger) error {

	dialSchema, err := getDailsSchema(client, driverType)
	if err != nil {
		return err
	}

	dialsSchemaLoader := gojsonschema.NewStringLoader(dialSchema)
	for _, dial := range driverInstance.Dials {
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

			logger.Error("driver-init", err, errData)

			return err
		}
	}

	return nil
}

func validateConfigSchema(client *rpc.Client, driverType string, configuration *json.RawMessage,
	logger lager.Logger) error {

	configSchema, err := getConfigSchema(client, driverType)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewStringLoader(configSchema)
	configLoader := gojsonschema.NewGoLoader(configuration)

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
		logger.Error("driver-init", err, errData)

		return err
	}

	return nil
}

func createClientAndCall(serviceMethod string, driverPath string, args interface{}, reply interface{}) error {
	client, err := createProviderClient(driverPath)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Call(serviceMethod, args, reply)
}

func createClientAndInvoke(call func(*rpc.Client) (string, error), driverPath string) (string, error) {
	client, err := createProviderClient(driverPath)
	if err != nil {
		return "", err
	}
	defer client.Close()

	return call(client)
}

func getDailsSchema(client *rpc.Client, driverType string) (string, error) {
	var result string
	err := client.Call(fmt.Sprintf("%s.GetDailsSchema", driverType), "", &result)

	return result, err
}

func getConfigSchema(client *rpc.Client, driverType string) (string, error) {
	var result string
	err := client.Call(fmt.Sprintf("%s.GetConfigSchema", driverType), "", &result)

	return result, err
}

func getDriverPath(driverType string) string {
	driverPath := os.Getenv("USB_DRIVER_PATH")
	if driverPath == "" {
		driverPath = "drivers"
	}

	driverPath = filepath.Join(driverPath, driverType)
	if runtime.GOOS == "windows" {
		driverPath = driverPath + ".exe"
	}

	return driverPath
}
