package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf/brokerapi"
)

type fileConfig struct {
	path   string
	loaded bool
	config *Config
}

func NewFileConfig(path string) ConfigProvider {
	return &fileConfig{path: path, loaded: false}
}

func (c *fileConfig) LoadConfiguration() (*Config, error) {
	var config *Config

	jsonConf, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}

	config, err = parseJson(jsonConf)
	if err != nil {
		return nil, err
	}
	c.loaded = true
	c.config = config
	if os.Getenv("PORT") != "" {
		c.config.BrokerAPI.Listen = ":" + os.Getenv("PORT")
	}
	return config, nil
}

func (c *fileConfig) GetDriverInstanceConfig(instanceID string) (*DriverInstance, error) {

	var instance DriverInstance
	exists := false

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, err
		}
	}

	for _, d := range c.config.Drivers {
		for _, di := range d.DriverInstances {
			if di.ID == instanceID {
				instance = *di
				exists = true
				break
			}
		}
	}

	if !exists {
		return nil, errors.New(fmt.Sprintf("Cannot find instanceID : %s", instanceID))
	}

	return &instance, nil

}

func (c *fileConfig) GetUaaAuthConfig() (*UaaAuth, error) {
	conf := (*json.RawMessage)(c.config.ManagementAPI.Authentication)

	uaa := Uaa{}
	err := json.Unmarshal(*conf, &uaa)
	if err != nil {
		return nil, err
	}
	return &uaa.UaaAuth, nil
}

func (c *fileConfig) GetDriver(driverID string) (Driver, error) {
	return Driver{}, nil
}

func (c *fileConfig) GetDriverInstance(instanceID string) (DriverInstance, error) {
	return DriverInstance{}, nil
}

func (c *fileConfig) GetService(instanceID string) (brokerapi.Service, error) {
	return brokerapi.Service{}, nil
}

func (c *fileConfig) GetDial(instanceID string, dialID string) (Dial, error) {
	return Dial{}, nil
}

func (c *fileConfig) SetDriver(driverInfo Driver) error {
	return nil
}

func (c *fileConfig) SetDriverInstance(driverID string, instance DriverInstance) error {
	return nil
}

func (c *fileConfig) SetService(instanceID string, service brokerapi.Service) error {
	return nil
}

func (c *fileConfig) SetDial(instanceID string, dialInfo Dial) error {
	return nil
}

func parseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
