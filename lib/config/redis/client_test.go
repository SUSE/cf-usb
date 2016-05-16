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
	err = redisConfig.testProvisioner.SetKV("usb", "{\"api_version\":\"2.6\",\"logLevel\":\"debug\",\"broker_api\":{\"external_url\":\"http://1.2.3.4:54054\",\"listen\":\":54054\",\"credentials\":{\"username\":\"username\",\"password\":\"password\"}},\"routes_register\":{\"nats_members\":[\"nats1\",\"nats2\"],\"broker_api_host\":\"broker\",\"management_api_host\":\"management\"},\"management_api\":{\"listen\":\":54053\",\"dev_mode\":false,\"broker_name\":\"usb\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"-----BEGIN PUBLIC KEY-----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB-----END PUBLIC KEY-----\"}},\"cloud_controller\":{\"api\":\"http://api.bosh-lite.com\",\"skip_tls_validation\":true}},\"instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000011\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}}", 5*time.Minute)
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
	value, err := redisConfig.testProvisioner.GetValue("usb")
	assert.NoError(err)
	t.Log(value)
}
