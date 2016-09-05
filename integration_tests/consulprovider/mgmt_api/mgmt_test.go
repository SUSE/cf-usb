package managementtest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"

	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"

	"github.com/hpcloud/cf-usb/lib/mgmt/operations"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger = lagertest.NewTestLogger("mgmt-api-test")
var sbMocked = new(sbMocks.USBServiceBroker)

var ConsulConfig = struct {
	ConsulAddress    string
	ConsulDatacenter string
	ConsulUser       string
	ConsulPassword   string
	ConsulSchema     string
	ConsulToken      string
}{}

var csmEndpoint = ""
var authToken = ""

func initiConsulProvisioner() (consul.Provisioner, error) {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")

	if ConsulConfig.ConsulAddress == "" {
		return nil, fmt.Errorf("CONSUL configuration environment variables not set")
	}

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
	return provisioner, err
}

func cleanupConsul() {
	provisioner, err := initiConsulProvisioner()
	if err != nil {
		logger.Error("Failed to cleanup", err)
	}
	err = provisioner.DeleteKVs("usb", nil)
	if err != nil {
		logger.Error("Failed to delete USB key", err)
	}
	logger.Info("cleanup_consul", lager.Data{"success": "finished cleaning consul"})

}

func initConsulProvider() (*config.Provider, error) {
	provisioner, err := initiConsulProvisioner()
	if err != nil {
		return nil, err
	}

	var list api.KVPairs
	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.6")})

	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}")})

	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}")})

	err = provisioner.PutKVs(&list, nil)
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

	csmClient := csm.NewCSMClient(logger)

	mgmt.ConfigureAPI(mgmtAPI, auth, provider, sbMocked, csmClient, logger, "t.t.t")

	return mgmtAPI, nil
}

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}
	defer cleanupConsul()
	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}
	response := mgmtInterface.GetInfoHandler.Handle(true)
	assert.IsType(&operations.GetInfoOK{}, response)
	info := response.(*operations.GetInfoOK).Payload

	assert.Equal("2.6", *info.BrokerAPIVersion)
}

func Test_RegisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}
	defer cleanupConsul()

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false, nil)
	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = uuid.NewV4().String()
	name := "testInstance"
	skiptls := true

	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = csmEndpoint
	params.DriverEndpoint.AuthenticationKey = authToken
	params.DriverEndpoint.SkipSSLValidation = &skiptls

	metadata := make(map[string]string)
	metadata["display_name"] = "servicename"

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
	t.Log(response)
}

func Test_UpdateInstanceEndpoint(t *testing.T) {
	assert := assert.New(t)
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}

	defer cleanupConsul()
	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}

	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}

	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}

	skiptls := true

	instanceID := uuid.NewV4().String()
	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = instanceID
	name := "testInstanceForUpdate"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = csmEndpoint
	params.DriverEndpoint.AuthenticationKey = authToken
	params.DriverEndpoint.SkipSSLValidation = &skiptls

	metadata := make(map[string]string)
	metadata["display_name"] = "servicename"

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
	paramsUpdate.DriverEndpoint.SkipSSLValidation = &skiptls
	params.DriverEndpoint.EndpointURL = csmEndpoint
	params.DriverEndpoint.AuthenticationKey = authToken
	params.DriverEndpoint.SkipSSLValidation = &skiptls

	responseUpdate := mgmtInterface.UpdateDriverEndpointHandler.Handle(*paramsUpdate, true)
	assert.IsType(&operations.UpdateDriverEndpointOK{}, responseUpdate)
}

func Test_GetDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}
	defer cleanupConsul()

	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}
	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}

	skiptls := true

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false, nil)
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
	params.DriverEndpoint.EndpointURL = csmEndpoint
	params.DriverEndpoint.AuthenticationKey = authToken
	params.DriverEndpoint.SkipSSLValidation = &skiptls

	metadata := make(map[string]string)
	metadata["display_name"] = "servicename"

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
	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}
	defer cleanupConsul()
	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}

	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}

	response := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, response)
}

func Test_UnregisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)

	consulProvider, err := initConsulProvider()
	if err != nil {
		t.Skip(err)
	}
	defer cleanupConsul()
	mgmtInterface, err := initMgmt(*consulProvider)
	if err != nil {
		t.Error(err)
	}

	if csmEndpoint == "" {
		t.Skip("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		t.Skip("CSM_API_KEY not set")
	}

	sbMocked.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)
	sbMocked.Mock.On("Delete", "usb").Return(nil)
	sbMocked.Mock.On("CheckServiceInstancesExist", mock.Anything).Return(false)

	skiptls := true

	instanceID := uuid.NewV4().String()
	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = instanceID
	name := "testUnregister"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = csmEndpoint
	params.DriverEndpoint.AuthenticationKey = authToken
	params.DriverEndpoint.SkipSSLValidation = &skiptls

	metadata := make(map[string]string)
	metadata["display_name"] = "servicename"

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
