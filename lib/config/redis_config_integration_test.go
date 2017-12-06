package config

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/SUSE/cf-usb/lib/brokermodel"
	"github.com/SUSE/cf-usb/lib/config/redis"

	"os"
	"strconv"
)

var RedisIntegrationConfig = struct {
	Provider Provider
	address  string
	password string
	db       int64
}{}

func init() {
	var err error
	RedisIntegrationConfig.address = os.Getenv("REDIS_ADDRESS")
	RedisIntegrationConfig.db = 0
	if os.Getenv("REDIS_DB") != "" {
		RedisIntegrationConfig.db, err = strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
		if err != nil {
			panic("REDIS_DB must be a 64bit integer")
		}
	}
}

func initRedisProvider() error {
	provisioner, err := redis.New(RedisIntegrationConfig.address, RedisIntegrationConfig.password, RedisIntegrationConfig.db)

	configSring, err := getRedisConfigString()
	if err != nil {
		return err
	}

	err = provisioner.SetKV("usb", configSring, 5*time.Minute)

	RedisIntegrationConfig.Provider = NewRedisConfig(provisioner)

	return err
}

func Test_RedisLoadConfiguration(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	config, err := RedisIntegrationConfig.Provider.LoadConfiguration()
	assert.NoError(err)

	assert.Equal("management", config.RoutesRegister.ManagmentAPIHost)
	assert.Equal("broker", config.RoutesRegister.BrokerAPIHost)

	assert.Equal(2, len(config.RoutesRegister.NatsMembers))
}

func Test_RedisGetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	instance, _, err := RedisIntegrationConfig.Provider.GetInstance("A0000000-0000-0000-0000-000000000004")
	assert.NoError(err)

	assert.Equal("local-mssql", instance.Name)
}

func Test_RedisGetDial(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	dial, instanceID, err := RedisIntegrationConfig.Provider.GetDial("C0000000-0000-0000-0000-000000000001")
	t.Log(instanceID)
	assert.NoError(err)

	assert.Equal("planmssql", dial.Plan.Name)
	assert.Equal("888B59E0-C2A2-4AB6-9335-2E90114A8F01", dial.Plan.ID)
}

func Test_RedisGetService(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	service, instanceID, err := RedisIntegrationConfig.Provider.GetService("83E94C97-C755-46A5-8653-461517EB442A")
	assert.NoError(err)

	assert.Equal("echo", service.Name)
	assert.Equal(true, service.Bindable)
	assert.Equal("echo Service", service.Description)
	assert.Equal("A0000000-0000-0000-0000-000000000002", instanceID)
}

func Test_RedisSetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var instance Instance
	instance.Name = "testDriverInstance"
	instance.Dials = make(map[string]Dial)
	instance.Service = brokermodel.CatalogService{}

	err = RedisIntegrationConfig.Provider.SetInstance("I0000000-0000-0000-0000-0000000000T1", instance)
	assert.NoError(err)
}

func Test_RedisSetService(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var service = brokermodel.CatalogService{}
	service.ID = "S0000000-0000-0000-0000-0000000000T1"
	service.Name = "testService"
	service.Tags = []string{"test"}
	service.Bindable = true
	service.Description = "test Service"

	var plan = brokermodel.Plan{}

	plan.ID = "P0000000-0000-0000-0000-0000000000T1"
	plan.Name = "free"
	plan.Free = true
	plan.Description = " test plan"
	plan.Metadata = &brokermodel.PlanMetadata{Metadata: struct{ DisplayName string }{"TestService"}} //.ServicePlanMetadata{DisplayName: "Test Service"}

	err = RedisIntegrationConfig.Provider.SetService("A0000000-0000-0000-0000-000000000002", service)
	assert.NoError(err)
}

func Test_RedisSetDial(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var plan brokermodel.Plan

	plan.ID = "P0000000-0000-0000-0000-0000000000T1"
	plan.Name = "free"
	plan.Free = true
	plan.Description = " test plan"
	plan.Metadata = &brokermodel.PlanMetadata{Metadata: struct{ DisplayName string }{"Test Service"}}

	var dial Dial
	dial.Plan = plan
	raw := json.RawMessage("{\"d1\":\"d2\"}")
	dial.Configuration = &raw

	err = RedisIntegrationConfig.Provider.SetDial("A0000000-0000-0000-0000-000000000002", "P0000000-0000-0000-0000-0000000000T1", dial)
	assert.NoError(err)
}
