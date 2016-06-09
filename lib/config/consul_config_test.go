package config

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/hpcloud/cf-usb/lib/brokermodel"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/consul/api"
	consulMock "github.com/hpcloud/cf-usb/lib/config/consul/mocks"
	"github.com/stretchr/testify/mock"
)

var TestConfig = struct {
	Provider Provider
}{}

func Test_SetDriverInstance(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var instance Instance
	instance.Name = "testInstance"
	err := TestConfig.Provider.SetInstance("testInstanceID", instance)
	assert.NoError(err)
}

func Test_GetDriverInstance(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var qoptions *api.QueryOptions

	provisioner.On("GetValue", "/usb/instances/testInstanceID/authentication_key").Return([]byte("secret_key"), nil)

	provisioner.On("GetValue", "/usb/instances/testInstanceID/target_url").Return([]byte("http://testurl.com:1234"), nil)

	provisioner.On("GetValue", "/usb/instances/testInstanceID/name").Return([]byte("MockedTestInstanceData"), nil)
	provisioner.On("GetValue", "/usb/instances/testInstanceID/service").Return([]byte("{\"id\":\"a\",\"Name\":\"test\"}"), nil)

	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	instance, parent, err := TestConfig.Provider.GetInstance("testInstanceID")
	t.Log(parent)
	assert.Equal("MockedTestInstanceData", instance.Name)
	assert.NoError(err)
}

func Test_DeleteDriverInstance(t *testing.T) {
	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner := new(consulMock.Provisioner)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)
	provisioner.On("DeleteKVs", "/usb/instances/testInstanceID", options).Return(nil)
	assert := assert.New(t)

	TestConfig.Provider = NewConsulConfig(provisioner)
	err := TestConfig.Provider.DeleteInstance("testInstance")
	assert.NoError(err)
}

func Test_SetService(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var serv = brokermodel.CatalogService{}
	serv.Bindable = true
	serv.Description = "testService desc"
	serv.ID = "testServiceID"
	serv.Metadata = &brokermodel.MetaData{DisplayName: "test service"} //struct{DisplayName string}{"test service"} //&brokerapi.ServiceMetadata{DisplayName: "test service"}
	serv.Name = "testService"
	serv.Tags = []string{"serv1", "serv2"}

	var plan = brokermodel.Plan{}
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokermodel.PlanMetadata{Metadata: struct{ DisplayName string }{"test plan"}}

	serv.Plans = []*brokermodel.Plan{&plan}

	err := TestConfig.Provider.SetService("testInstanceID", serv)
	assert.NoError(err)
}

func Test_GetService(t *testing.T) {
	var qoptions *api.QueryOptions
	provisioner := new(consulMock.Provisioner)

	provisioner.On("GetValue", "usb/instances/testInstanceID/authentication_key").Return([]byte("secret_key"), nil)

	provisioner.On("GetValue", "usb/api_version").Return([]byte("1.2"), nil)

	provisioner.On("GetValue", "usb/drivers_path").Return([]byte(""), nil)

	provisioner.On("GetValue", "usb/broker_api").Return([]byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}"), nil)

	provisioner.On("GetValue", "usb/management_api").Return([]byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}"), nil)

	provisioner.On("GetValue", "usb/instances/testInstanceID/target_url").Return([]byte("http://testurl.com:1234"), nil)

	provisioner.On("GetValue", "usb/instances/testInstanceID/name").Return([]byte("testInstance"), nil)

	provisioner.On("GetValue", "usb/instances/testInstanceID/service").Return([]byte("{\"Name\":\"testService\",\"id\":\"testServiceID\"}"), nil)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{
		"usb/instances/testInstanceID/name"}, nil)
	provisioner.On("GetAllKeys", "usb/instances/", "/", qoptions).Return([]string{
		"usb/instances/testInstanceID/"}, nil)
	provisioner.On("GetAllKeys", "usb/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)
	service, instanceID, err := TestConfig.Provider.GetService("testServiceID")
	t.Log(instanceID)
	t.Log(service)
	assert.Equal(service.Name, "testService")
	assert.NoError(err)
}

func Test_DeleteService(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)
	provisioner.On("DeleteKV", "/usb/instances/testInstanceID/service", options).Return(nil)
	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	err := TestConfig.Provider.DeleteService("testInstanceID")

	assert.NoError(err)
}

func Test_SetDial(t *testing.T) {

	provisioner := new(consulMock.Provisioner)

	var options *api.WriteOptions
	var qoptions *api.QueryOptions
	provisioner.On("AddKV", mock.Anything, mock.Anything, options).Return(nil)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{"/usb/instances/testInstanceID/name"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	var dialInfo Dial

	var plan brokermodel.Plan
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokermodel.PlanMetadata{Metadata: struct{ DisplayName string }{"test plan"}}

	raw := json.RawMessage("{\"a1\":\"b1\"}")

	dialInfo.Configuration = &raw
	dialInfo.Plan = plan

	err := TestConfig.Provider.SetDial("testInstanceID", "testDialID", dialInfo)
	assert.NoError(err)
}

func Test_GetDial(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var qoptions *api.QueryOptions
	provisioner.On("GetValue", mock.Anything).Return([]byte("{\"test\":\"dial\"}"), nil)
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/instances/testInstanceID/dials/dialID"}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	dialInfo, instanceID, err := TestConfig.Provider.GetDial("dialID")
	t.Log(instanceID)
	assert.NotNil(dialInfo)
	assert.NoError(err)
}

func Test_DeleteDial(t *testing.T) {
	provisioner := new(consulMock.Provisioner)

	var qoptions *api.QueryOptions
	var options *api.WriteOptions
	provisioner.On("GetAllKeys", mock.Anything, mock.Anything, qoptions).Return([]string{"/usb/instances/testInstanceID/dials/dialID"}, nil)
	provisioner.On("DeleteKV", "/usb/instances/testInstanceID/dials/dialID", options).Return(nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	assert := assert.New(t)

	err := TestConfig.Provider.DeleteDial("dialID")
	assert.NoError(err)
}

func Test_ConsulLoadConfig(t *testing.T) {
	assert := assert.New(t)
	provisioner := new(consulMock.Provisioner)
	var qoptions *api.QueryOptions

	provisioner.On("GetValue", "usb/instances/testInstanceID/authentication_key").Return([]byte("secret_key"), nil)

	provisioner.On("GetValue", "usb/api_version").Return([]byte("2.1"), nil)

	provisioner.On("GetValue", "usb/drivers_path").Return([]byte(""), nil)

	provisioner.On("GetValue", "usb/broker_api").Return([]byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}"), nil)

	provisioner.On("GetValue", "usb/management_api").Return([]byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}"), nil)

	provisioner.On("GetValue", "usb/instances/testInstanceID/name").Return([]byte("testInstance"), nil)
	provisioner.On("GetValue", "usb/instances/testInstanceID/target_url").Return([]byte("http://testurl.com:1234"), nil)

	provisioner.On("GetValue", "usb/instances/testInstanceID/service").Return([]byte("{\"Name\":\"testService\",\"id\":\"testServiceID\"}"), nil)
	provisioner.On("GetAllKeys", "usb/instances/", "", qoptions).Return([]string{
		"usb/instances/testInstanceID/name"}, nil)
	provisioner.On("GetAllKeys", "usb/instances/", "/", qoptions).Return([]string{
		"usb/instances/testInstanceID/"}, nil)

	provisioner.On("GetAllKeys", "usb/instances/testInstanceID/dials/", "/", qoptions).Return([]string{}, nil)

	TestConfig.Provider = NewConsulConfig(provisioner)

	config, err := TestConfig.Provider.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	t.Log(config.BrokerAPI)
	t.Log(config.ManagementAPI)
	t.Log(config.Instances)
	t.Log(config.APIVersion)
	assert.NoError(err)
}
