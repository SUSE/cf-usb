package mgmt

import (
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	"github.com/hpcloud/cf-usb/lib/operations"
)

var IntegrationConfig = struct {
	Provider         config.ConfigProvider
	MgmtAPI          *operations.UsbMgmtAPI
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

	IntegrationConfig.Provider = config.NewConsulConfig(provisioner)
	return true, nil
}

func initManager() error {
	swaggerJSON, err := data.Asset("swagger-spec/api.json")
	if err != nil {
		return err
	}

	swaggerSpec, err := spec.New(swaggerJSON, "")
	if err != nil {
		return err
	}

	IntegrationConfig.MgmtAPI = operations.NewUsbMgmtAPI(swaggerSpec)

	auth, err := uaa.NewUaaAuth("", "", true)
	if err != nil {
		return err
	}

	_, err = initProvider()

	if err != nil {
		return err
	}

	ConfigureAPI(IntegrationConfig.MgmtAPI, auth, IntegrationConfig.Provider)
	return nil
}

func Test_IntCreateDriver(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	var driver genmodel.Driver

	driver.Name = "testDriver"
	driver.ID = "testDriverID"
	driver.DriverType = "testDriverType"

	params := operations.CreateDriverParams{}
	params.Driver = driver

	info, err := IntegrationConfig.MgmtAPI.CreateDriverHandler.Handle(params)
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}
	t.Log(info)
	assert.NoError(err)
}

func Test_IntCreateDriverInstance(t *testing.T) {

	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateDriverInstanceParams{}
	var instace genmodel.DriverInstance

	instace.ID = "testInstanceID"
	instace.Name = "testInstanceName"
	instace.DriverID = "testDriverID"
	params.DriverInstance = instace
	info, err := IntegrationConfig.MgmtAPI.CreateDriverInstanceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntCreateDial(t *testing.T) {

	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateDialParams{}
	var instace genmodel.Dial

	instace.ID = "testDialID"
	instace.DriverInstanceID = "testInstanceID"
	instace.Plan = "{\"Name\":\"testPlanName\"}"
	params.Dial = instace

	info, err := IntegrationConfig.MgmtAPI.CreateDialHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntCreateService(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateServiceParams{}
	var instace genmodel.Service

	instace.ID = "testServiceID"
	instace.Bindable = true
	instace.DriverInstanceID = "testInstanceID"
	instace.Name = "testService"
	instace.Tags = []string{"test", "test service"}
	instace.Metadata = make(map[string]interface{})
	instace.Metadata["guid"] = "testGuid"

	params.Service = instace

	info, err := IntegrationConfig.MgmtAPI.CreateServiceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntCreateServicePlan(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.CreateServicePlanParams{}
	var plan genmodel.Plan

	plan.DialID = "testDialID"
	plan.Description = "testDescription"
	plan.Free = true
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	params.Plan = plan

	info, err := IntegrationConfig.MgmtAPI.CreateServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntUpdateServicePlan(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateServicePlanParams{}
	var plan genmodel.Plan

	plan.DialID = "testDialID"
	plan.Description = "testDescription Updated"
	plan.Free = true
	plan.ID = "testPlanID"
	plan.Name = "testPlanUpdated"

	params.PlanID = "testPlanID"
	params.Plan = plan

	info, err := IntegrationConfig.MgmtAPI.UpdateServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntUpdateDriver(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDriverParams{}

	var driver genmodel.Driver

	driver.Name = "testDriverUpdate"
	driver.ID = "testDriverID"
	driver.DriverType = "testDriverUpdateType"

	params.DriverID = driver.ID
	params.Driver = driver

	info, err := IntegrationConfig.MgmtAPI.UpdateDriverHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntUpdateDriverInstance(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDriverParams{}

	var driver genmodel.Driver

	driver.Name = "testDriverUpdate"
	driver.ID = "testDriverID"
	driver.DriverType = "testDriverUpdateType"

	params.DriverID = driver.ID
	params.Driver = driver

	info, err := IntegrationConfig.MgmtAPI.UpdateDriverHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntUpdateDial(t *testing.T) {

	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateDialParams{}

	var dial genmodel.Dial
	dial.DriverInstanceID = "testInstanceID"
	dial.ID = "testDialID"

	dial.Configuration = make(map[string]interface{})
	dial.Configuration["test"] = "test"

	params.Dial = dial
	params.DialID = dial.ID

	info, err := IntegrationConfig.MgmtAPI.UpdateDialHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

func Test_IntUpdateService(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.UpdateServiceParams{}
	var instace genmodel.Service

	instace.ID = "testServiceID"
	instace.Bindable = true
	instace.DriverInstanceID = "testInstanceID"
	instace.Name = "testUpdatedService"
	instace.Tags = []string{"test update", "test service"}
	instace.Metadata = make(map[string]interface{})
	instace.Metadata["guid"] = "testUpdateGuid"

	params.Service = instace
	params.ServiceID = instace.ID

	info, err := IntegrationConfig.MgmtAPI.UpdateServiceHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)
}

//Cleanup

func Test_IntDeleteDial(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.DeleteDialParams{}
	params.DialID = "testDialID"

	err = IntegrationConfig.MgmtAPI.DeleteDialHandler.Handle(params)
	assert.NoError(err)
}

func Test_IntDeleteServicePlan(t *testing.T) {

	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.DeleteServiceParams{}
	params.ServiceID = "testServiceID"
	err = IntegrationConfig.MgmtAPI.DeleteServiceHandler.Handle(params)
	assert.NoError(err)
}

func Test_IntDeleteDriverInstance(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.DeleteDriverInstanceParams{}
	params.DriverInstanceID = "testInstanceID"
	err = IntegrationConfig.MgmtAPI.DeleteDriverInstanceHandler.Handle(params)
	assert.NoError(err)
}

func Test_IntDeleteDriver(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt Create Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	params := operations.DeleteDriverParams{}
	params.DriverID = "testDriverID"
	err = IntegrationConfig.MgmtAPI.DeleteDriverHandler.Handle(params)
	assert.NoError(err)
}
