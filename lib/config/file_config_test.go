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
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type DummyServiceDials struct {
	MAXDB int `json:"max_dbsize_mb"`
}

func loadConfigAsset() (*Config, ConfigProvider, error) {
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
	assert.Equal(2, len(config.Drivers))
	assert.Equal("http://1.2.3.4:54054", config.BrokerAPI.ExternalUrl)
	assert.Equal(":54054", config.BrokerAPI.Listen)
	assert.Equal("username", config.BrokerAPI.Credentials.Username)
	assert.Equal("password", config.BrokerAPI.Credentials.Password)
	assert.Equal(":54053", config.ManagementAPI.Listen)
	assert.Equal(false, config.ManagementAPI.DevMode)
	assert.Equal("myuaaclient", config.ManagementAPI.UaaClient)
	assert.Equal("myuaasecret", config.ManagementAPI.UaaSecret)
	assert.Equal("http://api.bosh-lite.com", config.ManagementAPI.CloudController.Api)
	assert.Equal(true, config.ManagementAPI.CloudController.SkipTslValidation)
}

func TestLoadServiceConfig(t *testing.T) {
	assert := assert.New(t)

	config, _, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load config file")
	}

	for _, d := range config.Drivers {
		if d.DriverType == "dummy" {
			for _, di := range d.DriverInstances {
				dsp := DummyServiceProperties{}
				diConf := (*json.RawMessage)(di.Configuration)
				err := json.Unmarshal(*diConf, &dsp)
				if err != nil {
					assert.Equal(err, "Error unmarshaling properties")
				}
				if di.Name == "dummy1" {
					assert.Equal("one", dsp.PropOne)
					assert.Equal("two", dsp.PropTwo)
				}
				if di.Name == "dummy2" {
					assert.Equal("onenew", dsp.PropOne)
					assert.Equal("twonew", dsp.PropTwo)
				}
			}
		}
	}

}

func TestGetDriverConfig(t *testing.T) {
	assert := assert.New(t)

	_, provider, err := loadConfigAsset()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	driverInstance, err := provider.LoadDriverInstance("A0000000-0000-0000-0000-000000000002")
	if err != nil {
		assert.Error(err, "Unable to get driver configuration")
	}

	dsp := DummyServiceProperties{}

	conf := (*json.RawMessage)(driverInstance.Configuration)
	err = json.Unmarshal(*conf, &dsp)

	if err != nil {
		assert.Error(err, "Exception unmarshaling properties")
	}

	assert.Equal("one", dsp.PropOne)
	assert.Equal("two", dsp.PropTwo)

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

	assert.Equal(2, dials[0].MAXDB)
	assert.Equal(100, dials[1].MAXDB)
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
	assert.True(strings.Contains(uaaAuth.PublicKey, "public key"))
}
