package mgmt

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"

	"os"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/mock"
)

var testLogger *lagertest.TestLogger = lagertest.NewTestLogger("mgmt-api")

var IntegrationConfig = struct {
	Provider         config.ConfigProvider
	MgmtAPI          *operations.UsbMgmtAPI
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string

	CcServiceBroker *sbMocks.ServiceBrokerInterface
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

	workDir, err := os.Getwd()
	if err != nil {
		return false, err
	}
	buildDir := filepath.Join(workDir, "../../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

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

	IntegrationConfig.CcServiceBroker = new(sbMocks.ServiceBrokerInterface)

	ConfigureAPI(IntegrationConfig.MgmtAPI, auth, IntegrationConfig.Provider, IntegrationConfig.CcServiceBroker, testLogger)
	return nil
}

func Test_IntCreate(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt create test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	var driver genmodel.Driver

	driver.Name = "testDriver"
	driver.DriverType = "mysql"

	params := operations.CreateDriverParams{}
	params.Driver = driver

	info, err := IntegrationConfig.MgmtAPI.CreateDriverHandler.Handle(params)
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}
	t.Log(info)
	assert.NoError(err)

	IntegrationConfig.CcServiceBroker.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	instanceParams := operations.CreateDriverInstanceParams{}
	var instace genmodel.DriverInstance

	instace.Name = "testInstanceName"
	instace.DriverID = info.ID
	instanceConfig := make(map[string]interface{})

	instanceConfig["userid"] = "testUser"
	instanceConfig["password"] = "testPass"
	instanceConfig["server"] = "127.0.0.1"
	instanceConfig["port"] = "3306"

	instanceParams.DriverInstance = instace
	instanceParams.DriverInstance.DriverID = info.ID
	instanceParams.DriverInstance.Configuration = instanceConfig
	infoInstance, err := IntegrationConfig.MgmtAPI.CreateDriverInstanceHandler.Handle(instanceParams)
	t.Log(infoInstance)
	assert.NoError(err)

	dialParams := operations.CreateDialParams{}
	var instaceDial genmodel.Dial

	instaceDial.DriverInstanceID = infoInstance.ID
	instaceDial.Plan = "{\"Name\":\"testPlanName\"}"

	dialConfig := make(map[string]interface{})
	dialConfig["max_db_size_mb"] = "200"
	instaceDial.Configuration = dialConfig
	dialParams.Dial = instaceDial

	infoDial, err := IntegrationConfig.MgmtAPI.CreateDialHandler.Handle(dialParams)
	t.Log(infoDial)
	assert.NoError(err)

	serviceParams := operations.CreateServiceParams{}
	var instaceService genmodel.Service

	instaceService.Bindable = true
	instaceService.DriverInstanceID = infoInstance.ID
	instaceService.Name = "testService"
	instaceService.Tags = []string{"test", "test service"}
	instaceService.Metadata = make(map[string]interface{})

	serviceParams.Service = instaceService

	infoService, err := IntegrationConfig.MgmtAPI.CreateServiceHandler.Handle(serviceParams)
	t.Log(infoService)
	assert.NoError(err)

	splanParams := operations.CreateServicePlanParams{}
	var plan genmodel.Plan

	plan.DialID = infoDial.ID
	plan.Description = "testDescription"
	plan.Free = true
	plan.ID = "testPlanID"
	plan.Name = "testPlan"

	splanParams.Plan = plan

	infoPlan, err := IntegrationConfig.MgmtAPI.CreateServicePlanHandler.Handle(splanParams)
	t.Log(infoPlan)
	assert.NoError(err)
}

func Test_IntUpdate(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt update test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	callparams := operations.GetDriversParams{}

	drivers, err := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(callparams)
	if err != nil {
		t.Skip(err)
	}
	assert.NotNil(drivers)

	firstDriver := (*drivers)[0]

	t.Log(firstDriver)
	dialParams := operations.GetAllDialsParams{}
	dialParams.DriverInstanceID = firstDriver.DriverInstances[0]
	dials, err := IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(dialParams)
	if err != nil {
		t.Skip(err)
	}
	assert.NotNil(dials)

	firstDial := (*dials)[0]

	params := operations.UpdateServicePlanParams{}
	var plan genmodel.Plan

	plan.DialID = firstDial.ID
	plan.Description = "testDescription Updated"
	plan.Free = true
	plan.ID = firstDial.Plan
	plan.Name = "testPlanUpdated"

	params.PlanID = firstDial.Plan
	params.Plan = plan

	info, err := IntegrationConfig.MgmtAPI.UpdateServicePlanHandler.Handle(params)
	t.Log(info)
	assert.NoError(err)

	instanceParams := operations.GetDriverInstanceParams{}
	instanceParams.DriverInstanceID = firstDriver.DriverInstances[0]

	existingInstace, err := IntegrationConfig.MgmtAPI.GetDriverInstanceHandler.Handle(instanceParams)
	if err != nil {
		t.Skip(err)
	}

	serviceParams := operations.UpdateServiceParams{}
	var instace genmodel.Service

	instace.ID = existingInstace.Service
	instace.Bindable = true
	instace.DriverInstanceID = firstDriver.DriverInstances[0]
	instace.Name = "testUpdatedService"
	instace.Tags = []string{"test update", "test service"}
	instace.Metadata = make(map[string]interface{})
	instace.Metadata["guid"] = "testUpdateGuid"

	serviceParams.Service = instace
	serviceParams.ServiceID = instace.ID

	infoServiceUpdate, err := IntegrationConfig.MgmtAPI.UpdateServiceHandler.Handle(serviceParams)
	t.Log(infoServiceUpdate)
	assert.NoError(err)

	dialUpdateParams := operations.UpdateDialParams{}

	var dial genmodel.Dial
	dial.DriverInstanceID = existingInstace.ID
	dial.ID = firstDial.ID

	dial.Configuration = make(map[string]interface{})
	dial.Configuration["max_dbsize_mb"] = "400"

	dialUpdateParams.Dial = dial
	dialUpdateParams.DialID = dial.ID

	infoDial, err := IntegrationConfig.MgmtAPI.UpdateDialHandler.Handle(dialUpdateParams)
	t.Log(infoDial)
	assert.NoError(err)
}

//Cleanup

func Test_IntDelete(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt delete test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	callparams := operations.GetDriversParams{}

	drivers, err := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(callparams)
	if err != nil {
		t.Skip(err)
	}
	assert.NotNil(drivers)

	firstDriver := (*drivers)[0]

	dialParams := operations.GetAllDialsParams{}
	dialParams.DriverInstanceID = firstDriver.DriverInstances[0]
	dials, err := IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(dialParams)

	if err != nil {
		t.Skip(err)
	}
	assert.NotNil(dials)

	firstDial := (*dials)[0]

	dialDeleteParams := operations.DeleteDialParams{}
	dialDeleteParams.DialID = firstDial.ID

	instanceParams := operations.GetDriverInstanceParams{}
	instanceParams.DriverInstanceID = firstDriver.DriverInstances[0]

	existingInstace, err := IntegrationConfig.MgmtAPI.GetDriverInstanceHandler.Handle(instanceParams)
	if err != nil {
		t.Skip(err)
	}

	err = IntegrationConfig.MgmtAPI.DeleteDialHandler.Handle(dialDeleteParams)
	assert.NoError(err)

	deleteServiceParams := operations.DeleteServiceParams{}
	deleteServiceParams.ServiceID = existingInstace.Service

	err = IntegrationConfig.MgmtAPI.DeleteServiceHandler.Handle(deleteServiceParams)
	assert.NoError(err)

	deleteInstanceParams := operations.DeleteDriverInstanceParams{}
	deleteInstanceParams.DriverInstanceID = existingInstace.ID
	err = IntegrationConfig.MgmtAPI.DeleteDriverInstanceHandler.Handle(deleteInstanceParams)
	assert.NoError(err)

	deleteDriverParams := operations.DeleteDriverParams{}
	deleteDriverParams.DriverID = firstDriver.ID
	err = IntegrationConfig.MgmtAPI.DeleteDriverHandler.Handle(deleteDriverParams)
	assert.NoError(err)
}
