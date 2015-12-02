package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hpcloud/cf-usb/lib/config/redis"
	"github.com/pivotal-cf/brokerapi"

	"os"
	"strconv"
	"time"
)

var RedisIntegrationConfig = struct {
	Provider ConfigProvider
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

	err = provisioner.SetKV("broker_api", "{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", 5*time.Minute)
	if err != nil {
		return err
	}
	err = provisioner.SetKV("management_api", "{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", 5*time.Minute)
	if err != nil {
		return err
	}
	err = provisioner.SetKV("log_level", "debug", 5*time.Minute)
	if err != nil {
		return err
	}
	err = provisioner.SetKV("drivers", "[{\"driver_type\":\"dummy\",\"id\":\"00000000-0000-0000-0000-000000000001\",\"driver_instances\":[{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":[{\"id\":\"B0000000-0000-0000-0000-000000000001\",\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807D\",\"description\":\"This is the first plan\",\"free\":true}},{\"id\":\"B0000000-0000-0000-0000-000000000002\",\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F0D\",\"description\":\"This is the secondary plan\",\"free\":false}}],\"service\":{\"id\":\"GUID\",\"bindable\":true,\"name\":\"testService\",\"description\":\"test Service\",\"tags\":[\"testService\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}]}]", 5*time.Minute)
	if err != nil {
		return err
	}

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
	t.Log(config.LogLevel)
	assert.NoError(err)
}

func Test_RedisGetDriver(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	driver, err := RedisIntegrationConfig.Provider.GetDriver("00000000-0000-0000-0000-000000000001")
	t.Log(driver)
	assert.NoError(err)
}

func Test_RedisGetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	instance, err := RedisIntegrationConfig.Provider.GetDriverInstance("A0000000-0000-0000-0000-000000000002")
	t.Log(instance)
	assert.NoError(err)
}

func Test_RedisGetDial(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	dial, err := RedisIntegrationConfig.Provider.GetDial("A0000000-0000-0000-0000-000000000002", "B0000000-0000-0000-0000-000000000001")
	t.Log(dial)
	assert.NoError(err)
}

func Test_RedisGetService(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)
	service, err := RedisIntegrationConfig.Provider.GetService("A0000000-0000-0000-0000-000000000002")
	t.Log(service)
	assert.NoError(err)
}

func Test_RedisSetDriver(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var driver Driver
	driver.ID = "00000000-0000-0000-0000-0000000000T1"
	driver.DriverType = "testDriver"
	err = RedisIntegrationConfig.Provider.SetDriver(driver)
	assert.NoError(err)
}

func Test_RedisSetDriverInstance(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var instance DriverInstance
	instance.ID = "I0000000-0000-0000-0000-0000000000T1"
	instance.Name = "testDriverInstance"
	raw := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &raw

	err = RedisIntegrationConfig.Provider.SetDriverInstance("00000000-0000-0000-0000-0000000000T1", instance)
	assert.NoError(err)
}

func Test_RedisSetService(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var service brokerapi.Service
	service.ID = "S0000000-0000-0000-0000-0000000000T1"
	service.Name = "testService"
	service.Tags = []string{"test"}
	service.Bindable = true
	service.Description = "test Service"

	var plan brokerapi.ServicePlan

	plan.ID = "P0000000-0000-0000-0000-0000000000T1"
	plan.Name = "free"
	plan.Free = true
	plan.Description = " test plan"
	plan.Metadata = &brokerapi.ServicePlanMetadata{DisplayName: "Test Service"}

	err = RedisIntegrationConfig.Provider.SetService("I0000000-0000-0000-0000-0000000000T1", service)
	assert.NoError(err)
}

func Test_RedisSetDial(t *testing.T) {
	assert := assert.New(t)
	if RedisIntegrationConfig.address == "" {
		t.Skip("Skipping load configuration test : REDIS_ADDRESS must be set")
	}
	err := initRedisProvider()
	assert.NoError(err)

	var plan brokerapi.ServicePlan

	plan.ID = "P0000000-0000-0000-0000-0000000000T1"
	plan.Name = "free"
	plan.Free = true
	plan.Description = " test plan"
	plan.Metadata = &brokerapi.ServicePlanMetadata{DisplayName: "Test Service"}

	var dial Dial
	dial.Plan = plan
	dial.ID = "P0000000-0000-0000-0000-0000000000T1"
	raw := json.RawMessage("{\"d1\":\"d2\"}")
	dial.Configuration = &raw

	err = RedisIntegrationConfig.Provider.SetDial("I0000000-0000-0000-0000-0000000000T1", dial)
	assert.NoError(err)
}
