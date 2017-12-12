package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/SUSE/cf-usb/lib/brokermodel"
	"github.com/stretchr/testify/assert"
)

var MysqlIntegrationConfig = struct {
	address  string
	password string
	username string
	db       string
	Provider Provider
}{}

func initMysql() (bool, error) {
	MysqlIntegrationConfig.address = os.Getenv("MYSQL_ADDRESS")
	MysqlIntegrationConfig.db = os.Getenv("MYSQL_DB")
	MysqlIntegrationConfig.username = os.Getenv("MYSQL_USER")
	MysqlIntegrationConfig.password = os.Getenv("MYSQL_PASSWORD")

	if MysqlIntegrationConfig.address == "" || MysqlIntegrationConfig.db == "" || MysqlIntegrationConfig.username == "" || MysqlIntegrationConfig.password == "" {
		return true, nil
	}

	provider, err := NewMysqlConfig(MysqlIntegrationConfig.address, MysqlIntegrationConfig.username, MysqlIntegrationConfig.password, MysqlIntegrationConfig.db, "")
	if err != nil {
		return true, err
	}
	if err = provider.InitializeConfiguration(); err != nil {
		return true, err
	}
	MysqlIntegrationConfig.Provider = provider
	return false, nil
}

func Test_MysqlLoadConfiguration(t *testing.T) {
	assert := assert.New(t)
	skip, err := initMysql()
	if err != nil {
		t.Error(err)
	}
	if skip {
		t.Skip("MYSQL test environment variables not set")
	}
	config, err := MysqlIntegrationConfig.Provider.LoadConfiguration()
	assert.NoError(err)
	assert.NotNil(config)
	t.Log(config)
}

func Test_MysqlUaaConfig(t *testing.T) {
	assert := assert.New(t)
	skip, err := initMysql()
	if err != nil {
		t.Error(err)
	}
	if skip {
		t.Skip("MYSQL test environment variables not set")
	}
	uaa, err := MysqlIntegrationConfig.Provider.GetUaaAuthConfig()
	assert.NoError(err)
	assert.NotNil(uaa)
	t.Log(uaa)
}

func Test_MysqlInstanceTest(t *testing.T) {
	assert := assert.New(t)
	skip, err := initMysql()
	if err != nil {
		t.Error(err)
	}
	if skip {
		t.Skip("MYSQL test environment variables not set")
	}
	instance := Instance{}
	instance.AuthenticationKey = "authkey"
	instance.Name = "testInstance"
	instance.SkipSsl = true
	instance.TargetURL = "testInstance.test.com"
	instance.CaCert = ""

	var dial Dial
	var plan brokermodel.Plan
	plan.Name = "testPlan"
	plan.Description = "test plan description"
	plan.Free = true
	plan.ID = "testPlanID"
	var meta brokermodel.PlanMetadata
	meta.Name = "testMeta"
	meta.Description = "test meta description"
	plan.Metadata = &meta
	dial.Plan = plan
	raw := json.RawMessage("{\"a1\":\"b1\"}")

	dial.Configuration = &raw

	instance.Dials = make(map[string]Dial)

	instance.Dials["testDialGuid"] = dial

	var service brokermodel.CatalogService
	service.Name = "testService"
	service.Bindable = true
	service.Description = "testDescription"
	service.ID = "testServiceGuid"
	service.PlanUpdateable = true
	service.Tags = []string{"tag1", "tag2"}

	instance.Service = service

	notexists, err := MysqlIntegrationConfig.Provider.InstanceNameExists("testInstance")
	assert.NoError(err)
	assert.False(notexists)

	err = MysqlIntegrationConfig.Provider.SetInstance("testInstanceGuid", instance)
	assert.NoError(err)

	loadedInstance, instanceID, err := MysqlIntegrationConfig.Provider.GetInstance("testInstanceGuid")
	assert.NoError(err)
	t.Log(loadedInstance)
	assert.Equal(loadedInstance.Name, instance.Name)
	assert.Equal(instanceID, "testInstanceGuid")

	fullLoad, err := MysqlIntegrationConfig.Provider.LoadDriverInstance("testInstanceGuid")
	assert.NoError(err)
	assert.Equal(fullLoad.Service.Name, service.Name)

	exists, err := MysqlIntegrationConfig.Provider.InstanceNameExists("testInstance")
	assert.NoError(err)
	assert.True(exists)

	err = MysqlIntegrationConfig.Provider.DeleteInstance("testInstanceGuid")
	assert.NoError(err)
}
