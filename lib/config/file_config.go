package config

import (
	"encoding/json"
	"io/ioutil"
)

type fileConfig struct {
	ConfigProvider
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

	config, err = ParseJson(jsonConf)
	if err != nil {
		return nil, err
	}

	c.loaded = true
	c.config = config
	return config, nil
}

func ParseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (c *fileConfig) GetDriverProperties(driverType string) (DriverProperties, error) {
	var driverProperties DriverProperties
	var driverConfig DriverConfig

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return driverProperties, err
		}
	}

	for _, dc := range c.config.DriverConfigs {
		if dc.DriverType == driverType {
			driverConfig = dc
			break
		}
	}
	driverProperties.DriverConfiguration = driverConfig.Configuration

	for _, serviceID := range driverConfig.ServiceIDs {
		for _, service := range c.config.ServiceCatalog {
			if service.ID == serviceID {
				driverProperties.Services = append(driverProperties.Services, service)
				break
			}
		}
	}

	return driverProperties, nil
}

func (c *fileConfig) GetDriverTypes() ([]string, error) {
	var driverTypes []string

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return driverTypes, err
		}
	}

	for _, dc := range c.config.DriverConfigs {
		driverTypes = append(driverTypes, dc.DriverType)
	}
	return driverTypes, nil
}