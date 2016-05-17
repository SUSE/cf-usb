package managementtest

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hpcloud/cf-usb/lib/config/redis"

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

var RedisConfig = struct {
	RedisAddress  string
	RedisDatabase string
	RedisPassword string
}{}

func init() {
	RedisConfig.RedisAddress = os.Getenv("REDIS_ADDRESS")
	RedisConfig.RedisDatabase = os.Getenv("REDIS_DATABASE")
	RedisConfig.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if RedisConfig.RedisDatabase == "" {
		RedisConfig.RedisDatabase = "0"
	}
}

func initRedisProvider() (*config.Provider, error) {
	db, err := strconv.ParseInt(RedisConfig.RedisDatabase, 10, 64)
	if err != nil {
		return nil, err
	}
	provisioner, err := redis.New(RedisConfig.RedisAddress, RedisConfig.RedisPassword, db)
	if err != nil {
		return nil, err
	}

	err = provisioner.SetKV("usb", "{\"api_version\":\"2.6\",\"logLevel\":\"debug\",\"broker_api\":{\"external_url\":\"http://1.2.3.4:54054\",\"listen\":\":54054\",\"credentials\":{\"username\":\"username\",\"password\":\"password\"}},\"routes_register\":{\"nats_members\":[\"nats1\",\"nats2\"],\"broker_api_host\":\"broker\",\"management_api_host\":\"management\"},\"management_api\":{\"listen\":\":54053\",\"dev_mode\":false,\"broker_name\":\"usb\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"-----BEGIN PUBLIC KEY-----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB-----END PUBLIC KEY-----\"}},\"cloud_controller\":{\"api\":\"http://api.bosh-lite.com\",\"skip_tls_validation\":true}},\"instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000011\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}}", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	configProvider := config.NewRedisConfig(provisioner)

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

	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
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
	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
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
	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
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
	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
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
	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
	if err != nil {
		t.Error(err)
	}

	response := mgmtInterface.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, response)
}

func Test_UnregisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	if RedisConfig.RedisAddress == "" {
		t.Skip("Skipping management redis integration test - REDIS environment variables not set")
	}

	redisProvider, err := initRedisProvider()
	if err != nil {
		t.Error(err)
	}
	mgmtInterface, err := initMgmt(*redisProvider)
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
