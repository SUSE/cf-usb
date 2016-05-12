package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyServiceProperties struct {
	PropOne json.RawMessage `json:"service"`
}

type DummyServiceDials struct {
	MAXDB int `json:"max_dbsize_mb"`
}

func loadConfigAsset() (*Config, Provider, error) {
	config := &Config{}

	workDir, err := os.Getwd()
	configFile := filepath.Join(workDir, "../../test-assets/file-config/config.json")

	fileConfig := NewFileConfig(configFile)

	config, err = fileConfig.LoadConfiguration()
	if err != nil {
		return nil, nil, err
	}

	return config, fileConfig, nil
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	config, _, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load config file")
	}

	assert.Equal("2.6", config.APIVersion)
	assert.Equal("http://1.2.3.4:54054", config.BrokerAPI.ExternalURL)
	assert.Equal(":54054", config.BrokerAPI.Listen)
	assert.Equal("username", config.BrokerAPI.Credentials.Username)
	assert.Equal("password", config.BrokerAPI.Credentials.Password)
	assert.Equal(":54053", config.ManagementAPI.Listen)
	assert.Equal(false, config.ManagementAPI.DevMode)
	assert.Equal("myuaaclient", config.ManagementAPI.UaaClient)
	assert.Equal("myuaasecret", config.ManagementAPI.UaaSecret)
	assert.Equal("http://api.bosh-lite.com", config.ManagementAPI.CloudController.API)
	assert.Equal(true, config.ManagementAPI.CloudController.SkipTLSValidation)
}

func TestLoadServiceConfig(t *testing.T) {
	assert := assert.New(t)

	config, _, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load config file")
	}

	for _, d := range config.Instances {

		if d.Name == "dummy1" {
			assert.Equal("echo", d.Service.Name)
			assert.Equal("This is the first plan", d.Dials["B0000000-0000-0000-0000-000000000001"].Plan.Description)
		}
		if d.Name == "dummy2" {
			assert.Equal("echo", d.Service.Name)
			assert.Equal("This is the secondary plan", d.Dials["B0000000-0000-0000-0000-000000000011"].Plan.Description)
		}

	}

}

func TestGetDriverDials(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	driverInstance, err := configuration.LoadDriverInstance("A0000000-0000-0000-0000-000000000002")
	if err != nil {
		assert.Error(err, "Unable to get driver configuration")
	}

	var dials []DummyServiceDials

	for _, dial := range driverInstance.Dials {
		var dialDetails DummyServiceDials
		err = json.Unmarshal(*dial.Configuration, &dialDetails)
		if err != nil {
			assert.Error(err, "Exception unmarshaling dial configuration")
		}
		dials = append(dials, dialDetails)
	}
	assert.NotNil(dials[0].MAXDB)
	assert.NotNil(dials[1].MAXDB)
}

func TestGetUaaAuthConfig(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	uaaAuth, err := configuration.GetUaaAuthConfig()
	if err != nil {
		assert.Error(err, "Unable to get uaa auth config")
	}

	assert.Equal("usb.management.admin", uaaAuth.Scope)
	assert.True(strings.Contains(uaaAuth.PublicKey, "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d\nKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX\nqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug\nspULZVNRxq7veq/fzwIDAQAB\n-----END PUBLIC KEY-----"))
}

func TestDriverInstanceNameExists(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	exist, err := configuration.InstanceNameExists("dummy1")
	if err != nil {
		assert.Error(err, "Unable to check driver instance name existance")
	}

	assert.True(exist)
}
