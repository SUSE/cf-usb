package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/hashicorp/consul/api"
	consulMock "github.com/hpcloud/cf-usb/lib/config/consul/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/frodenas/brokerapi"
)

var TestConfig = struct {
	Provider ConfigProvider
}{}

func Test_ConsulSetDriver(t *testing.T) {

	provisioner := new(consulMock.ConsulProvisionerInterface)

	k := []byte("testType")
	n := []byte("testName")
	var options *api.WriteOptions
	provisioner.On("AddKV", "usb/drivers/testTypeID/Type", k, options).Return(nil)

	provisioner.On("AddKV", "usb/drivers/testTypeID/Name", n, options).Return(nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var driverInfo Driver
	driverInfo.DriverType = "testType"
	driverInfo.DriverName = "testName"

	err := TestConfig.Provider.SetDriver("testTypeID", driverInfo)
	assert.NoError(err)
}

func Test_GetDriver(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)
	var qoptions *api.QueryOptions

	k := []byte("testType")
	n := []byte("testName")
	i := []byte("testInstanceName")
	provisioner.On("GetValue", "usb/drivers/testID/Type").Return(k, nil)
	provisioner.On("GetValue", "usb/drivers/testID/Name").Return(n, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{
		"usb/drivers/testID/Name", "usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "/", qoptions).Return([]string{
		"usb/drivers/testID/"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/", "/", qoptions).Return([]string{
		"usb/drivers/testID/instances/testInstanceID/"}, nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Name").Return(i, nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Configuration").Return([]byte("{\"a\":\"b\"}"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/service").Return([]byte("{\"id\":\"a\",\"Name\":\"test\"}"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	driver, err := TestConfig.Provider.GetDriver("testID")
	assert.Equal("testType", string(driver.DriverType))
	assert.Equal("testName", string(driver.DriverName))
	assert.NoError(err)
}

func Test_DeleteDriver(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)
	var options *api.WriteOptions
	provisioner.On("DeleteKVs", "usb/drivers/testID", options).Return(nil)
	TestConfig.Provider = NewConsulConfig(provisioner)
	assert := assert.New(t)

	err := TestConfig.Provider.DeleteDriver("testID")
	assert.NoError(err)
}

func Test_SetDriverInstance(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var instance DriverInstance
	instance.Name = "testInstance"
	raw := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &raw
	err := TestConfig.Provider.SetDriverInstance("testID", "testInstanceID", instance)
	assert.NoError(err)
}

func Test_GetDriverInstance(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", "/usb/drivers/testID/instances/testInstanceID/Name").Return([]byte("MockedTestInstanceData"), nil)
	provisioner.On("GetValue", "/usb/drivers/testID/instances/testInstanceID/Configuration").Return([]byte("{\"test\":\"a\"}"), nil)
	provisioner.On("GetValue", "/usb/drivers/testID/instances/testInstanceID/service").Return([]byte("{\"id\":\"a\",\"Name\":\"test\"}"), nil)

	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	instance, parent, err := TestConfig.Provider.GetDriverInstance("testInstanceID")
	t.Log(parent)
	assert.Equal("MockedTestInstanceData", instance.Name)
	assert.NoError(err)
}

func Test_DeleteDriverInstance(t *testing.T) {
	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner := new(consulMock.ConsulProvisionerInterface)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("DeleteKVs", "/usb/drivers/testID/instances/testInstanceID", options).Return(nil)
	assert := assert.New(t)

	TestConfig.Provider = NewConsulConfig(provisioner)
	err := TestConfig.Provider.DeleteDriverInstance("testInstance")
	assert.NoError(err)
}

func Test_SetService(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var serv brokerapi.Service
	serv.Bindable = true
	serv.Description = "testService desc"
	serv.ID = "testServiceID"
	serv.Metadata = &brokerapi.ServiceMetadata{DisplayName: "test service"}
	serv.Name = "testService"
	serv.Tags = []string{"serv1", "serv2"}

	var plan brokerapi.ServicePlan
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokerapi.ServicePlanMetadata{DisplayName: "test plan"}

	serv.Plans = []brokerapi.ServicePlan{plan}

	err := TestConfig.Provider.SetService("testInstanceID", serv)
	assert.NoError(err)
}

func Test_GetService(t *testing.T) {
	var qoptions *api.QueryOptions
	provisioner := new(consulMock.ConsulProvisionerInterface)
	provisioner.On("GetValue", "usb/api_version").Return([]byte("1.2"), nil)

	provisioner.On("GetValue", "usb/drivers_path").Return([]byte(""), nil)

	provisioner.On("GetValue", "usb/broker_api").Return([]byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}"), nil)

	provisioner.On("GetValue", "usb/management_api").Return([]byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/Type").Return([]byte("testType"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/Name").Return([]byte("testName"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Name").Return([]byte("testInstance"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Configuration").Return([]byte("{\"a1\":\"b1\"}"), nil)
	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/service").Return([]byte("{\"Name\":\"testService\",\"id\":\"testServiceID\"}"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{
		"usb/drivers/testID/Name", "usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "/", qoptions).Return([]string{
		"usb/drivers/testID/"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/", "/", qoptions).Return([]string{
		"usb/drivers/testID/instances/testInstanceID/"}, nil)

	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)
	service, instanceID, err := TestConfig.Provider.GetService("testServiceID")
	t.Log(instanceID)
	t.Log(service)
	assert.Equal(service.Name, "testService")
	assert.NoError(err)
}

func Test_DeleteService(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("DeleteKV", "/usb/drivers/testID/instances/testInstanceID/service", options).Return(nil)
	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	err := TestConfig.Provider.DeleteService("testInstanceID")

	assert.NoError(err)
}

func Test_SetDial(t *testing.T) {

	provisioner := new(consulMock.ConsulProvisionerInterface)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var dialInfo Dial

	var plan brokerapi.ServicePlan
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokerapi.ServicePlanMetadata{DisplayName: "test plan"}

	raw := json.RawMessage("{\"a1\":\"b1\"}")

	dialInfo.Configuration = &raw
	dialInfo.Plan = plan

	err := TestConfig.Provider.SetDial("testInstanceID", "testDialID", dialInfo)
	assert.NoError(err)
}

func Test_GetDial(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", mock.Anything).Return([]byte("{\"test\":\"dial\"}"), nil)
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/dials/dialID"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	dialInfo, instanceID, err := TestConfig.Provider.GetDial("dialID")
	t.Log(instanceID)
	assert.NotNil(dialInfo)
	assert.NoError(err)
}

func Test_DeleteDial(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/dials/dialID"}, nil)
	provisioner.On("DeleteKV", "/usb/drivers/testID/instances/testInstanceID/dials/dialID", options).Return(nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	err := TestConfig.Provider.DeleteDial("dialID")
	assert.NoError(err)
}

func Test_ConsulLoadConfig(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(consulMock.ConsulProvisionerInterface)
	var qoptions *api.QueryOptions

	provisioner.On("GetValue", "usb/api_version").Return([]byte("2.1"), nil)

	provisioner.On("GetValue", "usb/drivers_path").Return([]byte(""), nil)

	provisioner.On("GetValue", "usb/broker_api").Return([]byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}"), nil)

	provisioner.On("GetValue", "usb/management_api").Return([]byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/Type").Return([]byte("testType"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/Name").Return([]byte("testName"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Name").Return([]byte("testInstance"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Configuration").Return([]byte("{\"a1\":\"b1\"}"), nil)
	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/service").Return([]byte("{\"Name\":\"testService\",\"id\":\"testServiceID\"}"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{
		"usb/drivers/testID/Name", "usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "/", qoptions).Return([]string{
		"usb/drivers/testID/"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/", "/", qoptions).Return([]string{
		"usb/drivers/testID/instances/testInstanceID/"}, nil)

	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	config, err := TestConfig.Provider.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	t.Log(config.BrokerAPI)
	t.Log(config.ManagementAPI)
	t.Log(config.Drivers)
	t.Log(config.APIVersion)
	assert.NoError(err)
}

func Test_DriverExists(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)
	var qoptions *api.QueryOptions

	k := []byte("testType")
	n := []byte("testName")
	i := []byte("testInstanceName")
	provisioner.On("GetValue", "usb/drivers/testID/Type").Return(k, nil)
	provisioner.On("GetValue", "usb/drivers/testID/Name").Return(n, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{
		"usb/drivers/testID/Name", "usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "/", qoptions).Return([]string{
		"usb/drivers/testID/"}, nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/", "/", qoptions).Return([]string{
		"usb/drivers/testID/instances/testInstanceID/"}, nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Name").Return(i, nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/Configuration").Return([]byte("{\"a\":\"b\"}"), nil)

	provisioner.On("GetValue", "usb/drivers/testID/instances/testInstanceID/service").Return([]byte("{\"id\":\"a\",\"Name\":\"test\"}"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/testID/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	exists, err := TestConfig.Provider.DriverExists("testID")
	assert.True(exists)
	assert.NoError(err)
}
