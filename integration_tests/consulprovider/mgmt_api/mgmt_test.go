package managementtest

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"

	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"

	"github.com/hpcloud/cf-usb/lib/mgmt/operations"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger = lagertest.NewTestLogger("mgmt-api-test")
var sbMocked = new(sbMocks.ServiceBrokerInterface)

var ConsulConfig = struct {
	ConsulAddress    string
	ConsulDatacenter string
	ConsulUser       string
	ConsulPassword   string
	ConsulSchema     string
	ConsulToken      string
}{}

func init() {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
}

func initConsulProvider() (*config.Provider, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.ConsulAddress
	consulConfig.Datacenter = ConsulConfig.ConsulDatacenter

	var auth api.HttpBasicAuth
	auth.Username = ConsulConfig.ConsulUser
	auth.Password = ConsulConfig.ConsulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = ConsulConfig.ConsulSchema

	consulConfig.Token = ConsulConfig.ConsulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return nil, err
	}
	configProvider := config.NewConsulConfig(provisioner)
	return &configProvider, nil
}

func initMgmt(provider config.Provider) (*operations.UsbMgmtAPI, error) {

	swaggerSpec, err := loads.Analyzed(mgmt.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)

	uaaAuthConfig, err := provider.GetUaaAuthConfig()
	if err != nil {
		logger.Error("initializing-uaa-config-failed", err)
		return nil, err
	}

	auth, err := uaa.NewUaaAuth(
		uaaAuthConfig.PublicKey,
		uaaAuthConfig.SymmetricVerificationKey,
		uaaAuthConfig.Scope,
		true,
		logger)
	if err != nil {
		logger.Fatal("initializing-uaa-auth-failed", err)
		return nil, err
	}

	mgmt.ConfigureAPI(mgmtAPI, auth, provider, sbMocked, logger, "t.t.t")

	return mgmtAPI, nil
}

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}

	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	response := mgmtInterface.GetInfoHandler.Handle(true)
	assert.IsType(&operations.GetInfoOK{}, response)
	info := response.(*operations.GetInfoOK).Payload

	assert.Equal("2.6", *info.BrokerAPIVersion)
}

func Test_RegisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = uuid.NewV4().String()
	name := "testInstance"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	metadata := &genmodel.EndpointMetadata{}
	metadata.DisplayName = "servicename"

	params.DriverEndpoint.Metadata = metadata

	var instanceDriver config.Instance
	instanceDriver.Name = "testInstance"

	var testConfig config.Config
	testConfig.Instances = make(map[string]config.Instance)
	testConfig.Instances["testInstanceID"] = instanceDriver
	testConfig.ManagementAPI = &config.ManagementAPI{}
	testConfig.ManagementAPI.BrokerName = "usb"

	response := mgmtInterface.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)
}

func Test_UpdateInstanceEndpoint(t *testing.T) {
	assert := assert.New(t)
	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	instanceID := uuid.NewV4().String()
	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = instanceID
	name := "testInstanceForUpdate"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	metadata := &genmodel.EndpointMetadata{}
	metadata.DisplayName = "servicename"

	params.DriverEndpoint.Metadata = metadata

	response := mgmtInterface.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)

	responseDrivers := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, responseDrivers)

	list := responseDrivers.(*operations.GetDriverEndpointsOK).Payload

	var updateInstanceID string
	for _, instance := range list {
		if *instance.Name == "testInstanceForUpdate" {
			updateInstanceID = instance.ID
			break
		}
	}

	paramsUpdate := &operations.UpdateDriverEndpointParams{}
	paramsUpdate.DriverEndpoint = &genmodel.DriverEndpoint{}
	paramsUpdate.DriverEndpoint.ID = updateInstanceID
	paramsUpdate.DriverEndpointID = updateInstanceID
	name = "updateName"
	paramsUpdate.DriverEndpoint.Name = &name
	paramsUpdate.DriverEndpoint.EndpointURL = "http://127.0.0.1:8081"
	paramsUpdate.DriverEndpoint.AuthenticationKey = "authkey"

	responseUpdate := mgmtInterface.UpdateDriverEndpointHandler.Handle(*paramsUpdate, true)
	assert.IsType(&operations.UpdateDriverEndpointOK{}, responseUpdate)
}

func Test_GetDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	instanceID := uuid.NewV4().String()
	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = instanceID
	name := "testGetInstance"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	metadata := &genmodel.EndpointMetadata{}
	metadata.DisplayName = "servicename"

	params.DriverEndpoint.Metadata = metadata

	response := mgmtInterface.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)

	responseDrivers := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, responseDrivers)

	list := responseDrivers.(*operations.GetDriverEndpointsOK).Payload

	var getInstanceID string
	for _, instance := range list {
		if *instance.Name == "testGetInstance" {
			getInstanceID = instance.ID
			break
		}
	}

	paramsGetDriver := &operations.GetDriverEndpointParams{}
	paramsGetDriver.DriverEndpointID = getInstanceID

	responseGet := mgmtInterface.GetDriverEndpointHandler.Handle(*paramsGetDriver, true)
	assert.IsType(&operations.GetDriverEndpointOK{}, responseGet)
}

func Test_GetDriverEndpoints(t *testing.T) {
	assert := assert.New(t)
	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}

	response := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, response)
}

func Test_UnregisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping management consul integration test - CONSUL environment variables not set")
	}
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Error(err)
	}

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)
	sbMocked.Mock.On("Delete", "usb").Return(nil)
	sbMocked.Mock.On("CheckServiceInstancesExist", mock.Anything).Return(false)

	instanceID := uuid.NewV4().String()
	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = instanceID
	name := "testUnregister"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	metadata := &genmodel.EndpointMetadata{}
	metadata.DisplayName = "servicename"

	params.DriverEndpoint.Metadata = metadata

	response := mgmtInterface.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)

	responseDrivers := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, responseDrivers)

	list := responseDrivers.(*operations.GetDriverEndpointsOK).Payload

	var getInstanceID string
	for _, instance := range list {
		if *instance.Name == "testUnregister" {
			getInstanceID = instance.ID
			break
		}
	}

	paramsUnregister := &operations.UnregisterDriverInstanceParams{}
	paramsUnregister.DriverEndpointID = getInstanceID

	responseUnregister := mgmtInterface.UnregisterDriverInstanceHandler.Handle(*paramsUnregister, true)
	assert.IsType(&operations.UnregisterDriverInstanceNoContent{}, responseUnregister)
}
