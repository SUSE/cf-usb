package mgmt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/frodenas/brokerapi"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/mocks"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"
	"github.com/hpcloud/cf-usb/lib/operations"
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
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	buildDir := filepath.Join(workDir, "../../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

	swaggerJSON, err := data.Asset("swagger-spec/api.json")
	if err != nil {
		return err
	}

	swaggerSpec, err := spec.New(swaggerJSON, "")
	if err != nil {
		return err
	}
	mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)

	auth, err := uaa.NewUaaAuth("", "", "", true, logger)
	if err != nil {
		return err
	}

	ConfigureAPI(mgmtAPI, auth, provider, sbMocked, logger)

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

	assert.Equal("2.6", info.Version)
}

func Test_CreateDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.CreateDriverParams{}

	driverId := "testDriverID"

	params.Driver = &genmodel.Driver{}
	params.Driver.ID = &driverId
	params.Driver.Name = "testDriver"
	params.Driver.DriverType = "testType"
	provider.On("DriverTypeExists", mock.Anything).Return(false, nil)
	provider.On("SetDriver", mock.Anything, mock.Anything).Return(nil)

	response := UnitTest.MgmtAPI.CreateDriverHandler.Handle(*params, true)
	assert.IsType(&operations.CreateDriverCreated{}, response)
}

func Test_CreateDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.CreateDriverInstanceParams{}
	params.DriverInstance = &genmodel.DriverInstance{}
	params.DriverInstance.DriverID = "testDriverID"
	testInstanceID := "testInstanceID"
	params.DriverInstance.ID = &testInstanceID
	params.DriverInstance.Name = "testInstance"
	params.DriverInstance.Configuration = map[string]interface{}{"property_one": "one", "property_two": "two"}

	var driver config.Driver
	var instanceDriver config.DriverInstance
	var dial config.Dial

	driver.DriverType = "dummy"

	dialConf := []byte(`{"configuration":{"max_dbsize_mb":2}}`)
	dial.Configuration = (*json.RawMessage)(&dialConf)

	instanceDriver.Name = "testInstance"
	instanceDriverConf := []byte(`{"property_one":"one", "property_two":"two"}`)
	instanceDriver.Configuration = (*json.RawMessage)(&instanceDriverConf)
	instanceDriver.Dials = make(map[string]config.Dial)
	instanceDriver.Dials["testDialID"] = dial

	driver.DriverInstances = make(map[string]config.DriverInstance)
	driver.DriverInstances["testInstanceID"] = instanceDriver

	provider.On("SetDriverInstance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	provider.On("SetDial", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	provider.On("GetDriver", "testDriverID").Return(&driver, nil)
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)
	provider.On("GetDriversPath").Return(os.Getenv("USB_DRIVER_PATH"), nil)

	var testConfig config.Config
	testConfig.Drivers = make(map[string]config.Driver)
	testConfig.Drivers["testDriverID"] = driver
	provider.On("LoadConfiguration").Return(&testConfig, nil)

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	sbMocked.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	response := UnitTest.MgmtAPI.CreateDriverInstanceHandler.Handle(*params, true)

	assert.IsType(&operations.CreateDriverInstanceCreated{}, response)
}

func Test_CreateDial(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.CreateDialParams{}
	params.Dial = &genmodel.Dial{}
	params.Dial.DriverInstanceID = "testInstanceID"
	dialID := "dialID"
	params.Dial.ID = &dialID
	planID := "planID"
	params.Dial.Plan = &planID

	provider.On("SetDial", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.CreateDialHandler.Handle(*params, true)
	assert.IsType(&operations.CreateDialCreated{}, response)
}

func Test_CreateServicePlan(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.CreateServicePlanParams{}
	params.Plan = &genmodel.Plan{}
	testDesc := "test desc"
	params.Plan.Description = &testDesc
	params.Plan.DialID = "testDialID"
	testPlanID := "testPlanID"
	params.Plan.ID = &testPlanID
	params.Plan.Name = "testPlan"
	pf := true
	params.Plan.Free = &pf

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	driver.DriverType = "testDriver"

	instace.Dials = make(map[string]config.Dial)
	instace.Dials["testDialID"] = dial
	driver.DriverInstances = make(map[string]config.DriverInstance)
	driver.DriverInstances["testInstanceID"] = instace
	testConfig.Drivers = make(map[string]config.Driver)
	testConfig.Drivers["testDriverID"] = driver

	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("DeleteDial", mock.Anything, mock.Anything).Return(nil)
	provider.On("SetDial", "testInstanceID", mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.CreateServicePlanHandler.Handle(*params, true)
	assert.IsType(&operations.CreateServicePlanCreated{}, response)
}

func Test_UpdateDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var driver config.Driver
	driver.DriverType = "testType"

	params := &operations.UpdateDriverParams{}
	params.Driver = &genmodel.Driver{}
	testDriverID := "testDriverID"
	params.Driver.ID = &testDriverID
	params.Driver.Name = "testDriver"
	params.Driver.DriverType = "testType"

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)

	provider.On("SetDriver", mock.Anything, mock.Anything).Return(nil)
	provider.On("GetDriver", mock.Anything).Return(&driver, nil)
	response := UnitTest.MgmtAPI.UpdateDriverHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateDriverOK{}, response)
}

func Test_UpdateDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	conf := json.RawMessage([]byte("{\"test\":\"test\"}"))
	var instanceInfo config.DriverInstance
	instanceInfo.Name = "testInstance"
	instanceInfo.Configuration = &conf
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID"}

	provider.On("GetDriverInstance", mock.Anything).Return(&instanceInfo, nil)

	params := &operations.UpdateDriverInstanceParams{}
	params.DriverConfig = &genmodel.DriverInstance{}
	params.DriverConfig.DriverID = "testDriverID"
	testInstanceID := "testInstanceID"
	params.DriverConfig.ID = &testInstanceID
	params.DriverConfig.Name = "testInstance"
	params.DriverInstanceID = "testDriverID"

	provider.On("SetDriverInstance", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.UpdateDriverInstanceHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateDriverInstanceOK{}, response)
}

func Test_UpdateDial(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.UpdateDialParams{}
	params.Dial = &genmodel.Dial{}
	params.DialID = "dialID"
	params.Dial.DriverInstanceID = "testInstanceID"
	updateddialID := "updateddialID"
	params.Dial.ID = &updateddialID
	planID := "planID"
	params.Dial.Plan = &planID

	var dial config.Dial
    provider.On("GetDial", updateddialID).Return(&dial, nil)
	provider.On("SetDial", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.UpdateDialHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateDialOK{}, response)
}

func Test_UpdateService(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	sbMocked.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)

	var testConfig config.Config
	var service brokerapi.Service
	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("GetService", "testServiceID").Return(&service, "testInstanceID", nil)

	params := &operations.UpdateServiceParams{}
	params.Service = &genmodel.Service{}
	params.ServiceID = "testServiceID"
	bindable := true
	params.Service.Bindable = &bindable
	params.Service.DriverInstanceID = "testInstanceID"
	updatedtestServiceID := "updatedtestServiceID"
	params.Service.ID = &updatedtestServiceID
	params.Service.Name = "updatedTestService"
	params.Service.Tags = []string{"test", "test Service"}
	desc := "description"
	params.Service.Description = &desc
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)
	provider.On("ServiceNameExists", mock.Anything).Return(false, nil)
	response := UnitTest.MgmtAPI.UpdateServiceHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateServiceOK{}, response)
}

func Test_UpdateServicePlan(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := &operations.UpdateServicePlanParams{}
	params.Plan = &genmodel.Plan{}
	params.PlanID = "testPlanID"
	testDesc := "test desc"
	params.Plan.Description = &testDesc
	params.Plan.DialID = "testDialID"
	testPlanID := "testPlanID"
	params.Plan.ID = &testPlanID
	params.Plan.Name = "testPlan"
	pf := true
	params.Plan.Free = &pf

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	driver.DriverType = "testDriver"
	instace.Dials = make(map[string]config.Dial)
	instace.Dials["testDialID"] = dial
	driver.DriverInstances = make(map[string]config.DriverInstance)
	driver.DriverInstances["testInstanceID"] = instace
	testConfig.Drivers = make(map[string]config.Driver)
	testConfig.Drivers["testDriverID"] = driver

	provider.On("GetPlan", params.PlanID).Return(&plan, mock.Anything, mock.Anything, nil)
	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("SetDial", "testInstanceID", mock.Anything, mock.Anything).Return(nil)
	response := UnitTest.MgmtAPI.UpdateServicePlanHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateServicePlanOK{}, response)
}

func Test_GetDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var instance config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	instance.Dials = make(map[string]config.Dial)
	instance.Dials["testDialID"] = dial

	var driverInfo config.Driver

	driverInfo.DriverType = "testDriverType"
	driverInfo.DriverInstances = make(map[string]config.DriverInstance)
	driverInfo.DriverInstances["testInstanceID"] = instance

	provider.On("GetDriver", mock.Anything).Return(&driverInfo, nil)

	params := &operations.GetDriverParams{}
	params.DriverID = "testDriverID"

	response := UnitTest.MgmtAPI.GetDriverHandler.Handle(*params, true)
	assert.IsType(&operations.GetDriverOK{}, response)
}

func Test_GetDrivers(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	driver.DriverType = "testDriver"
	instace.Dials = make(map[string]config.Dial)
	instace.Dials["testDialID"] = dial
	driver.DriverInstances = make(map[string]config.DriverInstance)
	driver.DriverInstances["testInstanceID"] = instace
	testConfig.Drivers = make(map[string]config.Driver)
	testConfig.Drivers["testDriverID"] = driver

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	response := UnitTest.MgmtAPI.GetDriversHandler.Handle(true)
	assert.IsType(&operations.GetDriversOK{}, response)
}

func Test_GetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	var instanceInfo config.DriverInstance
	instanceInfo.Name = "testInstance"
	instanceInfo.Configuration = &conf
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID"}

	provider.On("GetDriverInstance", mock.Anything).Return(&instanceInfo, nil)
	params := &operations.GetDriverInstanceParams{}
	params.DriverInstanceID = "testInstanceID"

	response := UnitTest.MgmtAPI.GetDriverInstanceHandler.Handle(*params, true)
	assert.IsType(&operations.GetDriverInstanceOK{}, response)
}

func Test_GetDriverInstances(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var firstInstance, secondInstance config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	firstInstance.Dials = make(map[string]config.Dial)
	firstInstance.Dials["testDialID"] = dial

	secondInstance.Dials = make(map[string]config.Dial)
	secondInstance.Dials["testDialID"] = dial

	var driverInfo config.Driver

	driverInfo.DriverType = "testDriverType"
	driverInfo.DriverInstances = make(map[string]config.DriverInstance)
	driverInfo.DriverInstances["testInstanceID"] = firstInstance
	driverInfo.DriverInstances["testSecondInstanceID"] = secondInstance

	provider.On("GetDriver", mock.Anything).Return(&driverInfo, nil)

	params := &operations.GetDriverInstancesParams{}
	params.DriverID = "testDriverID"

	response := UnitTest.MgmtAPI.GetDriverInstancesHandler.Handle(*params, true)
	assert.IsType(&operations.GetDriverInstancesOK{}, response)
}

func Test_GetDial(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	driver.DriverType = "testDriver"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf
	dial.Configuration = &conf

	instace.Dials = make(map[string]config.Dial)
	instace.Dials["testDialID"] = dial
	driver.DriverInstances = make(map[string]config.DriverInstance)
	driver.DriverInstances["testInstanceID"] = instace
	testConfig.Drivers = make(map[string]config.Driver)
	testConfig.Drivers["testDriverID"] = driver

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := &operations.GetDialParams{}
	params.DialID = "testDialID"

	response := UnitTest.MgmtAPI.GetDialHandler.Handle(*params, true)
	assert.IsType(&operations.GetDialOK{}, response)
}

func Test_GetAllDials(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	conf := json.RawMessage([]byte("{\"test\":\"test\"}"))

	var firstDial, secondDial config.Dial
	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	firstDial.Plan = plan
	secondDial.Plan = plan

	firstDial.Configuration = &conf
	secondDial.Configuration = &conf

	var instanceInfo config.DriverInstance
	instanceInfo.Name = "testInstance"
	instanceInfo.Configuration = &conf
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID"}
	instanceInfo.Dials = make(map[string]config.Dial)
	instanceInfo.Dials["testDialID"] = firstDial
	instanceInfo.Dials["testDialID2"] = secondDial

	provider.On("LoadDriverInstance", mock.Anything).Return(&instanceInfo, nil)

	params := &operations.GetAllDialsParams{}
	testInstanceID := "testInstanceID"
	params.DriverInstanceID = &testInstanceID

	response := UnitTest.MgmtAPI.GetAllDialsHandler.Handle(*params, true)
	assert.IsType(&operations.GetAllDialsOK{}, response)
}

func Test_GetService(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	var service brokerapi.Service

	service.ID = "testServiceID"
	service.Name = "testService"
	service.Plans = append(service.Plans, plan)

	provider.On("GetService", mock.Anything).Return(&service, "", nil)

	params := &operations.GetServiceParams{}
	params.ServiceID = "testServiceID"

	response := UnitTest.MgmtAPI.GetServiceHandler.Handle(*params, true)
	assert.IsType(&operations.GetServiceOK{}, response)
}

func Test_GetServiceByInstanceId(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	var instanceInfo config.DriverInstance
	instanceInfo.Name = "testInstance"
	instanceInfo.Configuration = &conf
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID", Name: "testService"}

	provider.On("GetDriverInstance", mock.Anything).Return(&instanceInfo, nil)
	params := &operations.GetServiceByInstanceIDParams{}
	params.DriverInstanceID = "testInstanceID"

	response := UnitTest.MgmtAPI.GetServiceByInstanceIDHandler.Handle(*params, true)
	assert.IsType(&operations.GetServiceByInstanceIDOK{}, response)
}

func Test_GetServicePlan(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	provider.On("GetPlan", mock.Anything).Return(&plan, "testDialId", "testInstanceID", nil)
	params := &operations.GetServicePlanParams{}
	params.PlanID = "testPlanID"

	response := UnitTest.MgmtAPI.GetServicePlanHandler.Handle(*params, true)
	assert.IsType(&operations.GetServicePlanOK{}, response)
}

func Test_GetServicePlans(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	conf := json.RawMessage([]byte("{\"test\":\"test\"}"))

	var firstDial, secondDial config.Dial
	var firstPlan, secondPlan brokerapi.ServicePlan

	firstPlan.ID = "testFirstPlanID"
	firstPlan.Name = "testFirstPlan"

	secondPlan.ID = "testSecondPlanID"
	secondPlan.Name = "testSecondPlan"

	firstDial.Plan = firstPlan
	secondDial.Plan = secondPlan

	firstDial.Configuration = &conf
	secondDial.Configuration = &conf

	var instanceInfo config.DriverInstance
	instanceInfo.Name = "testInstance"
	instanceInfo.Configuration = &conf
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID"}
	instanceInfo.Dials = make(map[string]config.Dial)
	instanceInfo.Dials["testDialID"] = firstDial
	instanceInfo.Dials["testDialID2"] = secondDial

	provider.On("LoadDriverInstance", mock.Anything).Return(&instanceInfo, nil)

	params := &operations.GetServicePlansParams{}
	testInstanceID := "testInstanceID"
	params.DriverInstanceID = &testInstanceID

	response := UnitTest.MgmtAPI.GetServicePlansHandler.Handle(*params, true)
	assert.IsType(&operations.GetServicePlansOK{}, response)
}
