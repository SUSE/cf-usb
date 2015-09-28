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
   "bolt_filename":"boltdb.db",
   "bolt_bucket":"brokerbucket",
   "api_version":"2.6",
   "auth_user":"demouser",
   "auth_password":"demopassword",
   "require_app_guid_in_bind_requests":true,
   "listen":":54054",
   "db_encryption_key":"12345678901234567890123456789012",
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
            "providerDisplayName":"MSSQL Service Ltd.",
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
            "providerDisplayName":"Echo Service Ltd.",
            "property_one":"one",
            "property_two":"two"
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
   ],
   "database":{  
      "parameters":{  
         "file":"boltdb.db",
         "bucket":"brokerbucket"
      }
   }
}
	`

type DummyServiceProperties struct {
	PropOne string `json:"property_one"`
	PropTwo string `json:"property_two"`
}

func writeTempConfigFile() (Config, error) {
	config := Config{}
	file, err := ioutil.TempFile(os.TempDir(), "loadconfig")
	if err != nil {
		return config, err
	}

	_, err = file.WriteString(configContent)
	if err != nil {
		return config, err
	}

	fileConfig := NewFileConfig(file.Name())

	config, err = fileConfig.LoadConfiguration()
	if err != nil {
		return config, err
	}
	defer os.Remove(file.Name())

	return config, nil
}

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	assert.Equal("2.6", config.APIVersion)
	assert.Equal(2, len(config.Services))
}

func TestLoadServiceConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := writeTempConfigFile()
	if err != nil {
		assert.Error(err, "Unable to load from temp config file")
	}

	for i := 0; i < len(config.Services); i++ {
		if config.Services[i].ID == "83E94C97-C755-46A5-8653-461517EB442A" {
			dsp := DummyServiceProperties{}
			err := json.Unmarshal(*config.Services[i].Metadata, &dsp)
			if err != nil {
				assert.Error(err, "Exception unmarshaling properties")
			}
			assert.Equal("one", dsp.PropOne)
			assert.Equal("two", dsp.PropTwo)
			break
		}
	}
}
