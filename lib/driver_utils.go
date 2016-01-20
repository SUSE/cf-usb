package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/natefinch/pie"
	"github.com/pivotal-golang/lager"
	"github.com/xeipuuv/gojsonschema"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"runtime"
)

func Validate(driverInstance config.DriverInstance, driverPath string, driverType string, logger lager.Logger) error {
	log := logger.Session("validate-driver-instance", lager.Data{"name": driverInstance.Name, "type": driverType})

	client, err := createProviderClient(getDriverPath(driverPath, driverType))
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

	pong, err := Ping(driverInstance.Configuration, driverPath, driverType)
	if err != nil {
		return err
	}

	if !pong {
		err = errors.New("Cannot reach server.")
		return err
	}

	return nil
}

func GetConfigSchema(driverPath string, driverType string, logger lager.Logger) (string, error) {
	log := logger.Session("get-driver-config-schema", lager.Data{"type": driverType})

	client, err := createProviderClient(getDriverPath(driverPath, driverType))
	if err != nil {
		return "", err
	}
	defer client.Close()

	log.Debug("get-config-schema", lager.Data{"type": driverType})

	schema, err := getConfigSchema(client, driverType)
	if err != nil {
		return "", err
	}
	return schema, nil
}

func GetDailsSchema(driverPath string, driverType string, logger lager.Logger) (string, error) {
	log := logger.Session("get-driver-dails-schema", lager.Data{"type": driverType})

	client, err := createProviderClient(getDriverPath(driverPath, driverType))
	if err != nil {
		return "", err
	}
	defer client.Close()

	log.Debug("get-dails-schema", lager.Data{"type": driverType})

	schema, err := getDailsSchema(client, driverType)
	if err != nil {
		return "", err
	}
	return schema, nil
}

func Ping(configuration *json.RawMessage, driverPath string, driverType string) (bool, error) {
	result := false
	path := getDriverPath(driverPath, driverType)
	err := createClientAndCall(fmt.Sprintf("%s.Ping", driverType), path, configuration, &result)

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

func getDriverPath(driverPath string, driverType string) string {

	driverPath = filepath.Join(driverPath, driverType)
	if runtime.GOOS == "windows" {
		driverPath = driverPath + ".exe"
	}

	return driverPath
}
