package mgmt

import (
	"os"
	"path/filepath"
	"testing"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/mocks"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"
	"github.com/hpcloud/cf-usb/lib/mgmt/operations"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mgmt-api-test")
var sbMocked *sbMocks.ServiceBrokerInterface = new(sbMocks.ServiceBrokerInterface)

var UnitTest = struct {
	MgmtAPI *operations.UsbMgmtAPI
}{}

func init_mgmt(provider config.ConfigProvider) error {

	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		return err
	}
	mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)

	auth, err := uaa.NewUaaAuth("", "", "", true, logger)
	if err != nil {
		return err
	}

	ConfigureAPI(mgmtAPI, auth, provider, sbMocked, logger, "t.t.t")

	UnitTest.MgmtAPI = mgmtAPI
	return nil
}

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	workDir, err := os.Getwd()
	configFile := filepath.Join(workDir, "../../test-assets/file-config/config.json")
	fileConfig := config.NewFileConfig(configFile)

	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	err = init_mgmt(fileConfig)
	if err != nil {
		t.Error(err)
	}
	response := UnitTest.MgmtAPI.GetInfoHandler.Handle(true)
	assert.IsType(&operations.GetInfoOK{}, response)
	info := response.(*operations.GetInfoOK).Payload

	assert.Equal("2.6", *info.BrokerAPIVersion)
	assert.Equal("t.t.t", *info.UsbVersion)
}

func Test_RegisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.RegisterDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = "testInstanceID"
	name := "testInstance"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	metadata := &genmodel.EndpointMetadata{}
	metadata.DisplayName = "servicename"

	params.DriverEndpoint.Metadata = metadata

	provider.On("SetDial", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)
	provider.On("SetInstance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	provider.On("InstanceNameExists", mock.Anything).Return(false, nil)

	var instanceDriver config.Instance
	instanceDriver.Name = "testInstance"

	var testConfig config.Config
	testConfig.Instances = make(map[string]config.Instance)
	testConfig.Instances["testInstanceID"] = instanceDriver
	testConfig.ManagementAPI = &config.ManagementAPI{}
	testConfig.ManagementAPI.BrokerName = "usb"
	provider.On("LoadConfiguration").Return(&testConfig, nil)

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	sbMocked.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	response := UnitTest.MgmtAPI.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)
}

func Test_UpdateInstanceEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var instanceInfo config.Instance
	instanceInfo.Name = "testInstance"

	var testConfig config.Config
	testConfig.Instances = make(map[string]config.Instance)
	testConfig.Instances["testInstanceID"] = instanceInfo
	testConfig.ManagementAPI = &config.ManagementAPI{}
	testConfig.ManagementAPI.BrokerName = "usb"
	provider.On("LoadConfiguration").Return(&testConfig, nil)

	provider.On("GetInstance", mock.Anything).Return(&instanceInfo, "testInstanceID", nil)

	params := &operations.UpdateDriverEndpointParams{}
	params.DriverEndpoint = &genmodel.DriverEndpoint{}
	params.DriverEndpoint.ID = "testInstanceID"
	name := "testInstance"
	params.DriverEndpoint.Name = &name
	params.DriverEndpoint.EndpointURL = "http://127.0.0.1:8080"
	params.DriverEndpoint.AuthenticationKey = "authkey"

	provider.On("InstanceNameExists", mock.Anything).Return(false, nil)
	provider.On("SetInstance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.UpdateDriverEndpointHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateDriverEndpointOK{}, response)
}

func Test_GetDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var instanceInfo config.Instance
	instanceInfo.Name = "testInstance"
	provider.On("GetInstance", mock.Anything).Return(&instanceInfo, "testInstanceID", nil)

	params := &operations.GetDriverEndpointParams{}
	params.DriverEndpointID = "testInstanceID"

	response := UnitTest.MgmtAPI.GetDriverEndpointHandler.Handle(*params, true)
	assert.IsType(&operations.GetDriverEndpointOK{}, response)
}

func Test_GetDriverEndpoints(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var firstInstance, secondInstance config.Instance
	firstInstance.Name = "instance1"
	secondInstance.Name = "instance2"

	var instances = make(map[string]config.Instance)
	instances["id1"] = firstInstance
	instances["id2"] = secondInstance

	var testConfig config.Config
	testConfig.Instances = instances
	testConfig.ManagementAPI = &config.ManagementAPI{}
	testConfig.ManagementAPI.BrokerName = "usb"
	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("GetInstance", mock.Anything).Return(&instances, nil)

	response := UnitTest.MgmtAPI.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, response)
}

func Test_UnregisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var testConfig config.Config
	testConfig.ManagementAPI = &config.ManagementAPI{}
	testConfig.ManagementAPI.BrokerName = "usb"
	provider.On("LoadConfiguration").Return(&testConfig, nil)

	var instanceInfo config.Instance
	instanceInfo.Name = "testInstance"
	provider.On("GetInstance", mock.Anything).Return(&instanceInfo, "testInstanceID", nil)
	provider.On("DeleteInstance", mock.Anything).Return(nil)

	sbMocked.Mock.On("Delete", testConfig.ManagementAPI.BrokerName).Return(nil)
	sbMocked.Mock.On("CheckServiceInstancesExist", mock.Anything).Return(false)

	params := &operations.UnregisterDriverInstanceParams{}
	params.DriverEndpointID = "testInstanceID"

	response := UnitTest.MgmtAPI.UnregisterDriverInstanceHandler.Handle(*params, true)
	assert.IsType(&operations.UnregisterDriverInstanceNoContent{}, response)
}
