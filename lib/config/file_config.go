package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func parseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
