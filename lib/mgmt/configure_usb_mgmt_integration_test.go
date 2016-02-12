package mgmt

import (
	"fmt"
	"path/filepath"
	"runtime"

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
	sbMocks "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi/mocks"
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

	auth, err := uaa.NewUaaAuth("", "", "", true, testLogger)
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
	params.Driver = &driver

	response := IntegrationConfig.MgmtAPI.CreateDriverHandler.Handle(params, true)

	assert.IsType(&operations.CreateDriverCreated{}, response)

	IntegrationConfig.CcServiceBroker.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)
	IntegrationConfig.CcServiceBroker.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)

	instanceParams := &operations.CreateDriverInstanceParams{}
	var instace genmodel.DriverInstance

	driverInstance := response.(*operations.CreateDriverCreated).Payload

	instace.Name = "testInstanceName"
	instace.DriverID = *driverInstance.ID
	instanceConfig := make(map[string]interface{})

	instanceConfig["userid"] = os.Getenv("MYSQL_USER")
	instanceConfig["password"] = os.Getenv("MYSQL_PASS")
	instanceConfig["server"] = os.Getenv("MYSQL_HOST")
	instanceConfig["port"] = os.Getenv("MYSQL_PORT")

	instanceParams.DriverInstance = &instace
	instanceParams.DriverInstance.DriverID = *driverInstance.ID
	instanceParams.DriverInstance.Configuration = instanceConfig
	response = IntegrationConfig.MgmtAPI.CreateDriverInstanceHandler.Handle(*instanceParams, true)
	assert.IsType(&operations.CreateDriverInstanceCreated{}, response)

	infoInstance := response.(*operations.CreateDriverInstanceCreated).Payload

	dialParams := &operations.CreateDialParams{}
	var instaceDial genmodel.Dial

	instaceDial.DriverInstanceID = *infoInstance.ID
	testplan := "{\"Name\":\"testPlanName\"}"
	instaceDial.Plan = &testplan

	dialConfig := make(map[string]interface{})
	dialConfig["max_db_size_mb"] = "200"
	instaceDial.Configuration = &dialConfig
	dialParams.Dial = &instaceDial

	response = IntegrationConfig.MgmtAPI.CreateDialHandler.Handle(*dialParams, true)
	assert.IsType(&operations.CreateDialCreated{}, response)

	infoDial := response.(*operations.CreateDialCreated).Payload
	t.Log(infoDial.Plan)
	//	serviceParams := operations.CreateServiceParams{}
	//	var instaceService genmodel.Service

	//	instaceService.Bindable = true
	//	instaceService.DriverInstanceID = infoInstance.ID
	//	instaceService.Name = "testService"
	//	instaceService.Tags = []string{"test", "test service"}
	//	instaceService.Metadata = make(map[string]interface{})

	//	serviceParams.Service = instaceService

	//	infoService, err := IntegrationConfig.MgmtAPI.CreateServiceHandler.Handle(serviceParams)
	//	t.Log(infoService)
	//	assert.NoError(err)
}

func Test_IntCreateDriverConflict(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt create test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	var driver, secondDriver genmodel.Driver
	driver.Name = "testDriver"
	driver.DriverType = "someType"

	params := operations.CreateDriverParams{}
	params.Driver = &driver

	response := IntegrationConfig.MgmtAPI.CreateDriverHandler.Handle(params, true)
	assert.IsType(&operations.CreateDriverCreated{}, response)

	createdDriver := response.(*operations.CreateDriverCreated).Payload

	secondDriver.Name = "someDriver"
	secondDriver.DriverType = "someType"

	params.Driver = &secondDriver

	response = IntegrationConfig.MgmtAPI.CreateDriverHandler.Handle(params, true)
	assert.IsType(&operations.CreateDriverConflict{}, response)

	deleteDriverParams := &operations.DeleteDriverParams{}
	deleteDriverParams.DriverID = *createdDriver.ID

	response = IntegrationConfig.MgmtAPI.DeleteDriverHandler.Handle(*deleteDriverParams, true)
	assert.IsType(&operations.DeleteDriverNoContent{}, response)
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

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)

	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload
	t.Log("Drivers === ", drivers)
	firstDriver := drivers[0]

	t.Log(firstDriver)
	dialParams := &operations.GetAllDialsParams{}

	for key, _ := range firstDriver.DriverInstances {
		dialParams.DriverInstanceID = &firstDriver.DriverInstances[key]
		break
	}

	response = IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(*dialParams, true)

	assert.IsType(&operations.GetAllDialsOK{}, response)

	dials := response.(*operations.GetAllDialsOK).Payload

	firstDial := dials[0]
	t.Log("FIRST DIAL:")
	t.Log(firstDial)

	params := &operations.UpdateServicePlanParams{}
	var plan genmodel.Plan

	plan.DialID = *firstDial.ID
	desc := "testDescription Updated"
	plan.Description = &desc
	pf := true
	plan.Free = &pf
	plan.ID = firstDial.Plan
	plan.Name = "testPlanUpdated"

	params.PlanID = *firstDial.Plan
	params.Plan = &plan
	t.Log("PLAN UPDATE")
	t.Log(plan)
	t.Log("PLAN ID", plan.ID)

	IntegrationConfig.CcServiceBroker.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)

	response = IntegrationConfig.MgmtAPI.UpdateServicePlanHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateServicePlanOK{}, response)

	instanceParams := &operations.GetDriverInstanceParams{}

	for key, _ := range firstDriver.DriverInstances {
		instanceParams.DriverInstanceID = firstDriver.DriverInstances[key]
		break
	}
	response = IntegrationConfig.MgmtAPI.GetDriverInstanceHandler.Handle(*instanceParams, true)

	assert.IsType(&operations.GetDriverInstanceOK{}, response)
	existingInstace := response.(*operations.GetDriverInstanceOK).Payload

	serviceParams := &operations.UpdateServiceParams{}
	var instace genmodel.Service

	instace.ID = existingInstace.Service
	bindable := true
	instace.Bindable = &bindable
	for key, _ := range firstDriver.DriverInstances {
		instace.DriverInstanceID = firstDriver.DriverInstances[key]
		break
	}
	instace.Name = "testUpdatedService"
	instace.Tags = []string{"test update", "test service"}
	description := "desc"
	instace.Description = &description
	instace.Metadata = make(map[string]interface{})
	instace.Metadata.(map[string]interface{})["guid"] = "testUpdateGuid"

	serviceParams.Service = &instace
	serviceParams.ServiceID = *instace.ID

	response = IntegrationConfig.MgmtAPI.UpdateServiceHandler.Handle(*serviceParams, true)
	assert.IsType(&operations.UpdateServiceOK{}, response)

	dialUpdateParams := &operations.UpdateDialParams{}

	var dial genmodel.Dial
	dial.DriverInstanceID = *existingInstace.ID
	dial.ID = firstDial.ID

	dial.Configuration = make(map[string]interface{})
	dial.Configuration.(map[string]interface{})["max_dbsize_mb"] = 400

	dial.Plan = plan.ID

	dialUpdateParams.Dial = &dial
	dialUpdateParams.DialID = *dial.ID

	response = IntegrationConfig.MgmtAPI.UpdateDialHandler.Handle(*dialUpdateParams, true)
	assert.IsType(&operations.UpdateDialOK{}, response)
}

func Test_IntServicePlan(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt update test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}
	IntegrationConfig.CcServiceBroker.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("EnableServiceAccess", mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("CheckServiceNameExists", mock.Anything).Return(false)

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)

	//get first driver
	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload

	firstDriver := drivers[0]

	t.Log(firstDriver)

	//get all dials
	dialParams := &operations.GetAllDialsParams{}

	for key, _ := range firstDriver.DriverInstances {
		dialParams.DriverInstanceID = &firstDriver.DriverInstances[key]
		break
	}

	response = IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(*dialParams, true)
	assert.IsType(&operations.GetAllDialsOK{}, response)

	dials := response.(*operations.GetAllDialsOK).Payload

	//create dial
	createDialParams := &operations.CreateDialParams{}
	var instaceDial genmodel.Dial

	instaceDial.DriverInstanceID = *dialParams.DriverInstanceID

	dialConfig := make(map[string]interface{})
	dialConfig["max_db_size_mb"] = "200"
	instaceDial.Configuration = &dialConfig
	createDialParams.Dial = &instaceDial

	response = IntegrationConfig.MgmtAPI.CreateDialHandler.Handle(*createDialParams, true)
	assert.IsType(&operations.CreateDialCreated{}, response)

	createdDial := response.(*operations.CreateDialCreated).Payload
	t.Log("CREATED DIAL", createdDial)

	//get all dials after create
	response = IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(*dialParams, true)
	assert.IsType(&operations.GetAllDialsOK{}, response)

	dialsAfter := response.(*operations.GetAllDialsOK).Payload
	assert.Len(dialsAfter, len(dials)+1, "The size of slice is not as expected")

	//get service plan
	getPlanParams := &operations.GetServicePlanParams{}
	getPlanParams.PlanID = *createdDial.Plan

	response = IntegrationConfig.MgmtAPI.GetServicePlanHandler.Handle(*getPlanParams, true)
	assert.IsType(&operations.GetServicePlanOK{}, response)

	existingPlan := response.(*operations.GetServicePlanOK).Payload
	assert.Equal(false, *existingPlan.Free, "Service plan should be false")
	assert.Equal(*createdDial.ID, existingPlan.DialID, "Dial Id for service plan should be "+*createdDial.ID)
	t.Log("EXISTING PLAN", existingPlan)

	params := &operations.UpdateServicePlanParams{}

	desc := "some updated plan description"
	existingPlan.Description = &desc
	existingPlan.Name = "someUpdatedPlan"

	params.PlanID = *existingPlan.ID
	params.Plan = existingPlan
	t.Log("PLAN UPDATE", existingPlan)

	response = IntegrationConfig.MgmtAPI.UpdateServicePlanHandler.Handle(*params, true)
	assert.IsType(&operations.UpdateServicePlanOK{}, response)
}

func Test_IntGetDriverSchema(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt delete test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)

	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload
	t.Log("Drivers === ", drivers)
	firstDriver := drivers[0]

	params := operations.GetDriverSchemaParams{}
	params.DriverID = *firstDriver.ID
	responseSchema := IntegrationConfig.MgmtAPI.GetDriverSchemaHandler.Handle(params, true)

	assert.IsType(&operations.GetDriverSchemaOK{}, responseSchema)
	schema := responseSchema.(*operations.GetDriverSchemaOK).Payload
	t.Log(schema)
}

func Test_IntGetDialSchema(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt delete test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)

	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload
	t.Log("Drivers === ", drivers)
	firstDriver := drivers[0]

	params := operations.GetDialSchemaParams{}
	params.DriverID = *firstDriver.ID
	responseSchema := IntegrationConfig.MgmtAPI.GetDialSchemaHandler.Handle(params, true)

	assert.IsType(&operations.GetDialSchemaOK{}, responseSchema)
	schema := responseSchema.(*operations.GetDialSchemaOK).Payload
	t.Log(schema)
}

func Test_IntPing(t *testing.T) {
	assert := assert.New(t)

	if IntegrationConfig.consulAddress == "" {
		t.Skip("Skipping mgmt delete test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
	}

	err := initManager()
	if err != nil {
		t.Error(err)
	}

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)

	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload
	t.Log("Drivers === ", drivers)
	firstDriver := drivers[0]

	instanceParams := &operations.GetDriverInstanceParams{}

	for key, _ := range firstDriver.DriverInstances {
		instanceParams.DriverInstanceID = firstDriver.DriverInstances[key]
		break
	}
	response = IntegrationConfig.MgmtAPI.GetDriverInstanceHandler.Handle(*instanceParams, true)

	assert.IsType(&operations.GetDriverInstanceOK{}, response)
	existingInstace := response.(*operations.GetDriverInstanceOK).Payload

	params := operations.PingDriverInstanceParams{}
	params.DriverInstanceID = *existingInstace.ID
	responseSchema := IntegrationConfig.MgmtAPI.PingDriverInstanceHandler.Handle(params, true)

	assert.IsType(&operations.PingDriverInstanceOK{}, responseSchema)
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

	response := IntegrationConfig.MgmtAPI.GetDriversHandler.Handle(true)
	assert.IsType(&operations.GetDriversOK{}, response)
	drivers := response.(*operations.GetDriversOK).Payload

	firstDriver := drivers[0]

	dialParams := &operations.GetAllDialsParams{}
	for key, _ := range firstDriver.DriverInstances {
		dialParams.DriverInstanceID = &firstDriver.DriverInstances[key]
		break
	}
	response = IntegrationConfig.MgmtAPI.GetAllDialsHandler.Handle(*dialParams, true)
	assert.IsType(&operations.GetAllDialsOK{}, response)
	dials := response.(*operations.GetAllDialsOK).Payload

	firstDial := dials[0]

	dialDeleteParams := &operations.DeleteDialParams{}
	dialDeleteParams.DialID = *firstDial.ID

	instanceParams := &operations.GetDriverInstanceParams{}
	for key, _ := range firstDriver.DriverInstances {

		instanceParams.DriverInstanceID = firstDriver.DriverInstances[key]
		break
	}
	response = IntegrationConfig.MgmtAPI.GetDriverInstanceHandler.Handle(*instanceParams, true)
	assert.IsType(&operations.GetDriverInstanceOK{}, response)
	existingInstace := response.(*operations.GetDriverInstanceOK).Payload

	response = IntegrationConfig.MgmtAPI.DeleteDialHandler.Handle(*dialDeleteParams, true)
	assert.IsType(&operations.DeleteDialNoContent{}, response)

	IntegrationConfig.CcServiceBroker.Mock.On("GetServiceBrokerGuidByName", mock.Anything).Return("aguid", nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Update", "aguid", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	IntegrationConfig.CcServiceBroker.Mock.On("Delete", mock.Anything).Return(nil)

	deleteInstanceParams := &operations.DeleteDriverInstanceParams{}
	deleteInstanceParams.DriverInstanceID = *existingInstace.ID
	response = IntegrationConfig.MgmtAPI.DeleteDriverInstanceHandler.Handle(*deleteInstanceParams, true)
	assert.IsType(&operations.DeleteDriverInstanceNoContent{}, response)

	deleteDriverParams := &operations.DeleteDriverParams{}
	deleteDriverParams.DriverID = *firstDriver.ID
	response = IntegrationConfig.MgmtAPI.DeleteDriverHandler.Handle(*deleteDriverParams, true)
	assert.IsType(&operations.DeleteDriverNoContent{}, response)
}
