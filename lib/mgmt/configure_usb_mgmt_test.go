package mgmt

import (
	"os"
	"path/filepath"
	"testing"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/mocks"
	//"github.com/hpcloud/cf-usb/lib/csm"
	csmMocks "github.com/hpcloud/cf-usb/lib/csm/mocks"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	//"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"
	"github.com/hpcloud/cf-usb/lib/mgmt/operations"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger = lagertest.NewTestLogger("mgmt-api-test")

type mockObjects struct {
	serviceBroker *sbMocks.USBServiceBroker
	csmClient     *csmMocks.CSM
	usbMgmt       *operations.UsbMgmtAPI
}

func initMgmt(provider config.Provider) (mockObjects, error) {

	mObjects := mockObjects{}
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		return mObjects, err
	}
	mObjects.usbMgmt = operations.NewUsbMgmtAPI(swaggerSpec)
	mObjects.csmClient = new(csmMocks.CSM)
	mObjects.serviceBroker = new(sbMocks.USBServiceBroker)

	auth, err := uaa.NewUaaAuth("", "", "", true, logger)
	if err != nil {
		return mObjects, err
	}

	ConfigureAPI(mObjects.usbMgmt, auth, provider, mObjects.serviceBroker, mObjects.csmClient, logger, "t.t.t")

	return mObjects, nil
}

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	workDir, err := os.Getwd()
	configFile := filepath.Join(workDir, "../../test-assets/file-config/config.json")
	fileConfig := config.NewFileConfig(configFile)

	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	mObjects, err := initMgmt(fileConfig)
	if err != nil {
		t.Error(err)
	}
	response := mObjects.usbMgmt.GetInfoHandler.Handle(true)
	assert.IsType(&operations.GetInfoOK{}, response)
	info := response.(*operations.GetInfoOK).Payload

	assert.Equal("2.6", *info.BrokerAPIVersion)
	assert.Equal("t.t.t", *info.UsbVersion)
}

func Test_RegisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.Provider)

	mObjects, err := initMgmt(provider)
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
	provider.On("SetDial", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)
	mObjects.serviceBroker.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	mObjects.serviceBroker.Mock.On("GetServiceBrokerGUIDByName", mock.Anything).Return("aguid", nil)
	mObjects.serviceBroker.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mObjects.serviceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mObjects.serviceBroker.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)
	mObjects.csmClient.Mock.On("Login", params.DriverEndpoint.EndpointURL, params.DriverEndpoint.AuthenticationKey, "", false).Return(nil)
	mObjects.csmClient.Mock.On("GetStatus").Return(nil)

	response := mObjects.usbMgmt.RegisterDriverEndpointHandler.Handle(*params, true)

	assert.IsType(&operations.RegisterDriverEndpointCreated{}, response)
}

func Test_UpdateInstanceEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.Provider)

	mObjects, err := initMgmt(provider)
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
	response := mObjects.usbMgmt.UpdateDriverEndpointHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateDriverEndpointOK{}, response)
}

func Test_GetDriverEndpoint(t *testing.T) {

	assert := assert.New(t)
	provider := new(mocks.Provider)

	mObjects, err := initMgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var instanceInfo config.Instance
	instanceInfo.Name = "testInstance"
	provider.On("GetInstance", mock.Anything).Return(&instanceInfo, "testInstanceID", nil)

	params := &operations.GetDriverEndpointParams{}
	params.DriverEndpointID = "testInstanceID"

	response := mObjects.usbMgmt.GetDriverEndpointHandler.Handle(*params, true)
	assert.IsType(&operations.GetDriverEndpointOK{}, response)

}

func Test_GetDriverEndpoints(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.Provider)

	mObjects, err := initMgmt(provider)
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

	response := mObjects.usbMgmt.GetDriverEndpointsHandler.Handle(true)
	assert.IsType(&operations.GetDriverEndpointsOK{}, response)
}

func Test_UnregisterDriverEndpoint(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.Provider)

	mObjects, err := initMgmt(provider)
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

	mObjects.serviceBroker.Mock.On("Delete", testConfig.ManagementAPI.BrokerName).Return(nil)
	mObjects.serviceBroker.Mock.On("CheckServiceInstancesExist", mock.Anything).Return(false)

	params := &operations.UnregisterDriverInstanceParams{}
	params.DriverEndpointID = "testInstanceID"

	response := mObjects.usbMgmt.UnregisterDriverInstanceHandler.Handle(*params, true)
	assert.IsType(&operations.UnregisterDriverInstanceNoContent{}, response)
}
