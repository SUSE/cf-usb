package mgmt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/mocks"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mgmt-api")
var sbMocked *sbMocks.ServiceBrokerInterface = new(sbMocks.ServiceBrokerInterface)

var UnitTest = struct {
	MgmtAPI *operations.UsbMgmtAPI
}{}

func init_mgmt(provider config.ConfigProvider) error {
	swaggerJSON, err := data.Asset("swagger-spec/api.json")
	if err != nil {
		return err
	}

	swaggerSpec, err := spec.New(swaggerJSON, "")
	if err != nil {
		return err
	}
	mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)

	auth, err := uaa.NewUaaAuth("", "", true)
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

	params := operations.GetInfoParams{""}

	info, err := UnitTest.MgmtAPI.GetInfoHandler.Handle(params)
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.Equal("2.6", info.Version)
}

func Test_CreateDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateDriverParams{}

	params.Driver.ID = "testDriverID"
	params.Driver.Name = "testDriver"
	params.Driver.DriverType = "testType"
	provider.On("SetDriver", mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.CreateDriverHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_CreateDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateDriverInstanceParams{}

	params.DriverInstance.DriverID = "testDriverID"
	params.DriverInstance.ID = "testInstanceID"
	params.DriverInstance.Name = "testInstance"

	var testDriver config.Driver
	testDriver.DriverType = "test"

	provider.On("SetDriverInstance", mock.Anything, mock.Anything).Return(nil)
	provider.On("SetDial", mock.Anything, mock.Anything).Return(nil)
	provider.On("GetDriver", "testDriverID").Return(testDriver, nil)
	_, err = UnitTest.MgmtAPI.CreateDriverInstanceHandler.Handle(params)
	t.Log("Expected error from validation call :", err)
	assert.Error(err)
}

func Test_CreateDial(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateDialParams{}

	params.Dial.DriverInstanceID = "testInstanceID"
	params.Dial.ID = "dialID"
	params.Dial.Plan = "planID"

	provider.On("SetDial", mock.Anything, mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.CreateDialHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_CreateService(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var testConfig config.Config

	params := operations.CreateServiceParams{}

	params.Service.Bindable = true
	params.Service.DriverInstanceID = "testInstanceID"
	params.Service.ID = "testServiceID"
	params.Service.Name = "testService"
	params.Service.Tags = []string{"test", "test Service"}
	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)

	sbMocked.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	sbMocked.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sbMocked.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	info, err := UnitTest.MgmtAPI.CreateServiceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_CreateServicePlan(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateServicePlanParams{}
	params.Plan.Description = "test desc"
	params.Plan.DialID = "testDialID"
	params.Plan.ID = "testPlanID"
	params.Plan.Name = "testPlan"

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"

	dial.ID = "testDialID"
	instace.ID = "testInstanceID"

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("DeleteDial", mock.Anything, mock.Anything).Return(nil)
	provider.On("SetDial", "testInstanceID", mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.CreateServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_UpdateDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDriverParams{}

	params.Driver.ID = "testDriverID"
	params.Driver.Name = "testDriver"
	params.Driver.DriverType = "testType"
	provider.On("SetDriver", mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.UpdateDriverHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_UpdateDriverInstance(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDriverInstanceParams{}

	params.DriverConfig.DriverID = "testDriverID"
	params.DriverConfig.ID = "testInstanceID"
	params.DriverConfig.Name = "testInstance"
	params.DriverInstanceID = "testDriverID"

	provider.On("SetDriverInstance", mock.Anything, mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.UpdateDriverInstanceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_UpdateDial(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDialParams{}
	params.DialID = "dialID"

	params.Dial.DriverInstanceID = "testInstanceID"
	params.Dial.ID = "updateddialID"
	params.Dial.Plan = "planID"

	provider.On("SetDial", mock.Anything, mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.UpdateDialHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_UpdateService(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateServiceParams{}
	params.ServiceID = "testServiceID"
	params.Service.Bindable = true
	params.Service.DriverInstanceID = "testInstanceID"
	params.Service.ID = "updatedtestServiceID"
	params.Service.Name = "updatedTestService"
	params.Service.Tags = []string{"test", "test Service"}
	provider.On("SetService", mock.Anything, mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.UpdateServiceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_UpdateServicePlan(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateServicePlanParams{}
	params.PlanID = "testPlanID"
	params.Plan.Description = "test desc"
	params.Plan.DialID = "testDialID"
	params.Plan.ID = "testPlanID"
	params.Plan.Name = "testPlan"

	var testConfig config.Config

	var driver config.Driver
	var instace config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)
	provider.On("SetDial", "testInstanceID", mock.Anything).Return(nil)
	info, err := UnitTest.MgmtAPI.UpdateServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetDriver(t *testing.T) {
	assert := assert.New(t)
	provider := new(mocks.ConfigProvider)

	err := init_mgmt(provider)
	if err != nil {
		t.Error(err)
	}

	var instace config.DriverInstance
	var dial config.Dial

	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan

	dial.ID = "testDialID"
	instace.ID = "testInstanceID"

	instace.Dials = append(instace.Dials, dial)

	var driverInfo config.Driver

	driverInfo.ID = "testDriverID"
	driverInfo.DriverType = "testDriverType"
	driverInfo.DriverInstances = append(driverInfo.DriverInstances, &instace)

	provider.On("GetDriver", mock.Anything).Return(driverInfo, nil)

	params := operations.GetDriverParams{}
	params.DriverID = "testDriverID"

	info, err := UnitTest.MgmtAPI.GetDriverHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
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

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetDriversParams{}

	info, err := UnitTest.MgmtAPI.GetDriversHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
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
	instanceInfo.ID = "testInstanceID"
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID"}

	provider.On("GetDriverInstance", mock.Anything).Return(instanceInfo, nil)
	params := operations.GetDriverInstanceParams{}
	params.DriverInstanceID = "testInstanceID"

	info, err := UnitTest.MgmtAPI.GetDriverInstanceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetDriverInstances(t *testing.T) {
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

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetDriverInstancesParams{}
	params.DriverID = "testDriverID"

	info, err := UnitTest.MgmtAPI.GetDriverInstancesHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
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

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf
	dial.Configuration = &conf

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetDialParams{}
	params.DialID = "testDialID"

	info, err := UnitTest.MgmtAPI.GetDialHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetAllDials(t *testing.T) {
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
	var dial2 config.Dial
	var plan brokerapi.ServicePlan
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	dial.Plan = plan
	dial2.Plan = plan

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	dial2.ID = "testDialID2"
	instace.ID = "testInstanceID"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf
	dial.Configuration = &conf
	dial2.Configuration = &conf

	instace.Dials = append(instace.Dials, dial)

	instace.Dials = append(instace.Dials, dial2)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetAllDialsParams{}
	params.DriverInstanceID = "testInstanceID"

	info, err := UnitTest.MgmtAPI.GetAllDialsHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetService(t *testing.T) {
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

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf
	dial.Configuration = &conf

	var service brokerapi.Service

	service.ID = "testServiceID"
	service.Name = "testService"
	service.Plans = append(service.Plans, plan)

	instace.Service = service

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetServiceParams{}
	params.ServiceID = "testServiceID"

	info, err := UnitTest.MgmtAPI.GetServiceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetServices(t *testing.T) {
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
	instanceInfo.ID = "testInstanceID"
	instanceInfo.Service = brokerapi.Service{ID: "testServiceID", Name: "testService"}

	provider.On("GetDriverInstance", mock.Anything).Return(instanceInfo, nil)
	params := operations.GetServicesParams{}
	params.DriverInstanceID = "testInstanceID"

	info, err := UnitTest.MgmtAPI.GetServicesHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_GetServicePlan(t *testing.T) {
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

	driver.ID = "testDriverID"
	driver.DriverType = "testDriver"
	dial.ID = "testDialID"
	instace.ID = "testInstanceID"
	var conf json.RawMessage

	conf = json.RawMessage([]byte("{\"test\":\"test\"}"))

	instace.Configuration = &conf
	dial.Configuration = &conf

	var service brokerapi.Service

	service.ID = "testServiceID"
	service.Name = "testService"
	service.Plans = append(service.Plans, plan)

	instace.Service = service

	instace.Dials = append(instace.Dials, dial)

	driver.DriverInstances = append(driver.DriverInstances, &instace)

	testConfig.Drivers = append(testConfig.Drivers, driver)

	provider.On("LoadConfiguration").Return(&testConfig, nil)

	params := operations.GetServicePlanParams{}
	params.PlanID = "testPlanID"

	info, err := UnitTest.MgmtAPI.GetServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}
