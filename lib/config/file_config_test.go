package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configContent = `
{  
   "api_version":"2.6",
	"broker_credentials": {
		"username": "username",
		"password": "password"
	},
	"logLevel": "debug",
   "require_app_guid_in_bind_requests":true,
   "listen":":54054",
	"management_listen":":54053",
	"start_mgmt": true,
   "db_encryption_key":"12345678901234567890123456789012",
   "driver_configs":[  
      {  
         "driver_type":"mssql",
         "configuration":{  
            "brokerGoSqlDriver":"mssql",
            "brokerMssqlConnection":{  
               "server":"127.0.0.1",
               "port":"38017",
               "database":"master",
               "user id":"sa",
               "password":"password1234!"
            },
            "servedMssqlBindingHostname":"192.168.1.10",
            "servedMssqlBindingPort":38017
         },
         "service_ids":[  
            "83E94C97-C755-46A5-8653-461517EB442C"
         ]
      },
      {  
         "driver_type":"dummy",
         "configuration":{  
            "property_one":"one",
            "property_two":"two"
         },
		 "dials": [
        {
            "planId": "53425178-F731-49E7-9E53-5CF4BE9D807A",
            "configuration": {
                "max_dbsize_mb": 2
            }
        },
        {
            "planId": "888B59E0-C2A1-4AB6-9335-2E90114A8F07",
            "configuration": {
                "max_dbsize_mb": 100
            }
        }
		],
         "service_ids":[  
            "83E94C97-C755-46A5-8653-461517EB442A"
         ]
      }
   ],
   "services":[  
      {  
         "id":"83E94C97-C755-46A5-8653-461517EB442C",
         "bindable":true,
         "name":"mssql",
         "description":"MSSQL Service",
         "tags":[  
            "mssql",
            "mssql"
         ],
         "metadata":{  
            "providerDisplayName":"MSSQL Service Ltd."
         },
         "plans":[  
            {  
               "name":"free",
               "id":"53425178-F731-49E7-9E53-5CF4BE9D807A",
               "description":"This is the first plan",
               "free":true
            },
            {  
               "name":"secondary",
               "id":"888B59E0-C2A1-4AB6-9335-2E90114A8F07",
               "description":"This is the secondary plan",
               "free":false
            }
         ]
      },
      {  
         "id":"83E94C97-C755-46A5-8653-461517EB442A",
         "bindable":true,
         "name":"echo",
         "description":"echo Service",
         "tags":[  
            "echo"
         ],
         "metadata":{  
            "providerDisplayName":"Echo Service Ltd."
         },
         "plans":[  
            {  
               "name":"free",
               "id":"53425178-F731-49E7-9E53-5CF4BE9D807A",
               "description":"This is the first plan",
               "free":true
            },
            {  
               "name":"secondary",
               "id":"888B59E0-C2A1-4AB6-9335-2E90114A8F07",
               "description":"This is the secondary plan",
               "free":false
            }
         ]
      }
   ]
}
	`

type DummyServiceProperties struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

type DummyServiceDials struct {
	MAXDB int `json:"max_dbsize_mb"`
}

func writeTempConfigFile() (*Config, ConfigProvider, error) {
	config := &Config{}
	file, err := ioutil.TempFile(os.TempDir(), "loadconfig")
	if err != nil {
		return config, nil, err
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(configContent)
	if err != nil {
		return config, nil, err
	}

	fileConfig := NewFileConfig(file.Name())

	config, err = fileConfig.LoadConfiguration()
	if err != nil {
		return config, nil, err
	}

	return config, fileConfig, nil
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	config, _, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	assert.Equal("2.6", config.APIVersion)
	assert.Equal(2, len(config.ServiceCatalog))
	assert.Equal(2, len(config.DriverConfigs))
	assert.Equal(":54054", config.Listen)
	assert.Equal(":54053", config.ManagementListen)
	assert.Equal("username", config.Crednetials.Username)
	assert.Equal("password", config.Crednetials.Password)
	assert.Equal("debug", config.LogLevel)
	assert.Equal(true, config.StartMgmt)
}

func TestLoadServiceConfig(t *testing.T) {
	assert := assert.New(t)

	config, _, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	for i := 0; i < len(config.DriverConfigs); i++ {
		if config.DriverConfigs[i].DriverType == "dummy" {
			dsp := DummyServiceProperties{}
			conf := (*json.RawMessage)(config.DriverConfigs[i].Configuration)
			err := json.Unmarshal(*conf, &dsp)
			if err != nil {
				assert.Error(err, "Exception unmarshaling properties")
			}
			assert.Equal("one", dsp.PropOne)
			assert.Equal("two", dsp.PropTwo)
			break
		}
	}
}

func TestGetDriverConfig(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	dummyProperties, err := configuration.GetDriverProperties("dummy")
	if err != nil {
		assert.Error(err, "Unable to get driver configuration")
	}

	dsp := DummyServiceProperties{}
	conf := (*json.RawMessage)(dummyProperties.DriverConfiguration)
	err = json.Unmarshal(*conf, &dsp)

	if err != nil {
		assert.Error(err, "Exception unmarshaling properties")
	}

	assert.Equal(1, len(dummyProperties.Services))
	assert.Equal("83E94C97-C755-46A5-8653-461517EB442A", dummyProperties.Services[0].ID)
	assert.Equal("one", dsp.PropOne)
	assert.Equal("two", dsp.PropTwo)

}

func TestGetDriverDials(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	dummyProperties, err := configuration.GetDriverProperties("dummy")
	if err != nil {
		assert.Error(err, "Unable to get driver configuration")
	}

	var dials []DummyServiceDials

	for _, dial := range dummyProperties.DriverDialsConfiguration {
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

func GetDriverTypes(t *testing.T) {
	assert := assert.New(t)

	_, configuration, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	driverTypes, err := configuration.GetDriverTypes()

	assert.Equal(2, len(driverTypes))
}
