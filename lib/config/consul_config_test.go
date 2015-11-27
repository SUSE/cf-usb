package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"

	"github.com/hashicorp/consul/api"
	consulMock "github.com/hpcloud/cf-usb/lib/config/consul/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/pivotal-cf/brokerapi"
)

var TestConfig = struct {
	Provider ConfigProvider
}{}

func Test_ConsulSetDriver(t *testing.T) {

	provisioner := new(consulMock.ConsulProvisionerInterface)

	k := []byte("testType")
	var options *api.WriteOptions
	provisioner.On("AddKV", "usb/drivers/testID", k, options).Return(nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var driverInfo Driver
	driverInfo.ID = "testID"
	driverInfo.DriverType = "testType"

	err := TestConfig.Provider.SetDriver(driverInfo)
	assert.NoError(err)
}

func Test_GetDriver(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	k := []byte("testType")
	provisioner.On("GetValue", "usb/drivers/testID").Return(k, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	driver, err := TestConfig.Provider.GetDriver("testID")
	assert.Equal("testType", string(driver.DriverType))
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
	instance.ID = "testInstanceID"
	instance.Name = "testInstance"
	raw := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &raw
	err := TestConfig.Provider.SetDriverInstance("testID", instance)
	assert.NoError(err)
}

func Test_GetDriverInstance(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", mock.Anything).Return([]byte("testInstance"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	instance, err := TestConfig.Provider.GetDriverInstance("testInstanceID")

	assert.Equal("testInstanceID", instance.ID)
	assert.Equal("testInstance", instance.Name)
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
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", mock.Anything).Return([]byte("{\"Name\":\"testService\"}"), nil)
	provisioner.On("GetAllKeys", "usb/drivers/", "", qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	service, err := TestConfig.Provider.GetService("testInstanceID")

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
	dialInfo.ID = "dialID"
	dialInfo.Plan = plan

	err := TestConfig.Provider.SetDial("testInstanceID", dialInfo)
	assert.NoError(err)
}

func Test_GetDial(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", mock.Anything).Return([]byte("{\"test\":\"dial\"}"), nil)
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	dialInfo, err := TestConfig.Provider.GetDial("testInstanceID", "dialID")
	assert.NotNil(dialInfo)
	assert.NoError(err)
}

func Test_DeleteDial(t *testing.T) {
	provisioner := new(consulMock.ConsulProvisionerInterface)

	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)
	provisioner.On("DeleteKV", "/usb/drivers/testID/instances/testInstanceID/dials/dialID", options).Return(nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	err := TestConfig.Provider.DeleteDial("testInstanceID", "dialID")
	assert.NoError(err)
}

func Test_ConsulLoadConfig(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(consulMock.ConsulProvisionerInterface)
	var qoptions *api.QueryOptions

	provisioner.On("GetValue", mock.Anything).Return([]byte("{\"test\":\"dial\"}"), nil)
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/drivers/testID/instances/testInstanceID/Name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	config, err := TestConfig.Provider.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	t.Log(config.BrokerAPI)
	t.Log(config.ManagementAPI)
	t.Log(config.Drivers)
	t.Log(config.LogLevel)

	assert.NoError(err)
}
