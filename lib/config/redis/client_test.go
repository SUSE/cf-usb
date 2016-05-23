package redis

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var redisConfig = struct {
	address         string
	password        string
	db              int64
	testProvisioner Provisioner
}{}

func init() {
	var err error
	redisConfig.address = os.Getenv("REDIS_ADDRESS")
	redisConfig.password = os.Getenv("REDIS_PASSWORD")
	redisConfig.db = 0
	if os.Getenv("REDIS_DB") != "" {
		redisConfig.db, err = strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
		if err != nil {
			panic("REDIS_DB must be a 64bit integer")
		}
	}
}

func initProvider() error {
	var err error
	redisConfig.testProvisioner, err = New(redisConfig.address, redisConfig.password, redisConfig.db)
	return err
}

func Test_SetValue(t *testing.T) {
	assert := assert.New(t)

	if redisConfig.address == "" {
		t.Skip("Skipping set value test : REDIS_ADDRESS must be set")
	}
	err := initProvider()
	if err != nil {
		t.Error(err)
	}
	err = redisConfig.testProvisioner.SetKV("broker_api", "{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", 5*time.Minute)
	err = redisConfig.testProvisioner.SetKV("management_api", "{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}", 5*time.Minute)
	err = redisConfig.testProvisioner.SetKV("drivers", "[{\"driver_type\":\"dummy\",\"id\":\"00000000-0000-0000-0000-000000000001\",\"driver_instances\":[{\"name\":\"dummy1\",\"id\":\"A0000000-0000-0000-0000-000000000002\",\"configuration\":{\"property_one\":\"one\",\"property_two\":\"two\"},\"dials\":[{\"id\":\"B0000000-0000-0000-0000-000000000001\",\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807D\",\"description\":\"This is the first plan\",\"free\":true}},{\"id\":\"B0000000-0000-0000-0000-000000000002\",\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F0D\",\"description\":\"This is the secondary plan\",\"free\":false}}],\"service\":{\"id\":\"GUID\",\"bindable\":true,\"name\":\"testService\",\"description\":\"test Service\",\"tags\":[\"testService\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}]}]", 5*time.Minute)
	assert.NoError(err)
}

func Test_GetValue(t *testing.T) {
	assert := assert.New(t)

	if redisConfig.address == "" {
		t.Skip("Skipping set value test : REDIS_ADDRESS must be set")
	}

	err := initProvider()
	if err != nil {
		t.Error(err)
	}
	value, err := redisConfig.testProvisioner.GetValue("broker_api")
	assert.NoError(err)
	t.Log(value)
}
