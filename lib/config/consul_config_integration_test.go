package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hpcloud/cf-usb/lib/config/consul"

	"github.com/frodenas/brokerapi"
	"github.com/hashicorp/consul/api"
	"os"
)

var IntegrationConfig = struct {
	Provider         ConfigProvider
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string
}{}

func init() {
	IntegrationConfig.consulAddress = os.Getenv("CONSUL_ADDRESS")
	IntegrationConfig.consulDatacenter = os.Getenv("CONSUL_DATACENTER")
	IntegrationConfig.consulPassword = os.Getenv("CONSUL_PASSWORD")
	IntegrationConfig.consulUser = os.Getenv("CONSUL_USER")
	IntegrationConfig.consulSchema = os.Getenv("CONSUL_SCHEMA")
	IntegrationConfig.consulToken = os.Getenv("CONSUL_TOKEN")
}

func initProvider() (bool, error) {
	var consulConfig api.Config
	if IntegrationConfig.consulAddress == "" {
		return false, nil
	}
	consulConfig.Address = IntegrationConfig.consulAddress
	consulConfig.Datacenter = IntegrationConfig.consulPassword

	var auth api.HttpBasicAuth
	auth.Username = IntegrationConfig.consulUser
	auth.Password = IntegrationConfig.consulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = IntegrationConfig.consulSchema

	consulConfig.Token = IntegrationConfig.consulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return false, err
	}

	IntegrationConfig.Provider = NewConsulConfig(provisioner)
	return true, nil
}

func Test_IntConsulSetDriver(t *testing.T) {
	initialized, err := initProvider()

	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	var driverInfo Driver
	driverInfo.DriverType = "testType"
	driverInfo.DriverName = "testName"
	err = IntegrationConfig.Provider.SetDriver("testID", driverInfo)
	assert.NoError(err)
}

func Test_IntGetDriver(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	driver, err := IntegrationConfig.Provider.GetDriver("testID")
	assert.Equal("testType", string(driver.DriverType))
	assert.NoError(err)
}

func Test_IntSetDriverInstance(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	var instance DriverInstance
	instance.Name = "testInstance"
	raw := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &raw
	err = IntegrationConfig.Provider.SetDriverInstance("testID", "testInstanceID", instance)
	assert.NoError(err)
}

func Test_IntGetDriverInstance(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	instance, _, err := IntegrationConfig.Provider.GetDriverInstance("testInstanceID")

	assert.Equal("testInstance", instance.Name)
	assert.NoError(err)
}

func Test_IntLoadDriverInstance(t *testing.T) {
	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	instance, err := IntegrationConfig.Provider.LoadDriverInstance("testInstanceID")
	t.Log("Load driver instance results:")
	t.Log(instance.Configuration)
	t.Log(instance.Dials)
	t.Log(instance.Service)
	assert.Equal("testInstance", instance.Name)
	assert.NoError(err)
}

func Test_IntSetService(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

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

	err = IntegrationConfig.Provider.SetService("testInstanceID", serv)
	assert.NoError(err)
}

func Test_IntGetService(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	service, instanceID, err := IntegrationConfig.Provider.GetService("testServiceID")

	assert.Equal(instanceID, "testInstanceID")
	assert.Equal(service.Name, "testService")
	assert.Equal(service.Plans[0].Name, "free")
	assert.NoError(err)
}

func Test_IntSetDial(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

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

	err = IntegrationConfig.Provider.SetDial("testInstanceID", "testdialID", dialInfo)
	assert.NoError(err)
}

func Test_IntGetDial(t *testing.T) {

	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	dialInfo, instanceID, err := IntegrationConfig.Provider.GetDial("testdialID")
	t.Log(dialInfo)
	t.Log(instanceID)
	assert.NoError(err)
}

func Test_IntDriverInstanceNameExists(t *testing.T) {
	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	exist, err := IntegrationConfig.Provider.DriverInstanceNameExists("testInstance")
	if err != nil {
		assert.Error(err, "Unable to check driver instance name existance")
	}
	assert.NoError(err)
	assert.True(exist)
}

func Test_IntDriverTypeExists(t *testing.T) {
	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	exist, err := IntegrationConfig.Provider.DriverTypeExists("testType")
	if err != nil {
		assert.Error(err, "Unable to check driver type existance")
	}
	assert.NoError(err)
	assert.True(exist)
}

func Test_IntLoadConfiguration(t *testing.T) {
	initialized, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	configuration, err := IntegrationConfig.Provider.LoadConfiguration()
	if err != nil {
		assert.Error(err, "Load configuration failed")
	}
	t.Log(configuration.Drivers)
	assert.NoError(err)
}
