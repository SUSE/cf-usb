package config

import "github.com/hpcloud/gocfbroker"

type fileConfig struct {
	ConfigProvider
	path   string
	loaded bool
	config Config
}

func NewFileConfig(path string) ConfigProvider {
	return &fileConfig{path: path, loaded: false}
}

func (c *fileConfig) LoadConfiguration() (Config, error) {
	var config Config

	err := gocfbroker.LoadConfig(c.path, &config)
	if err != nil {
		return config, err
	}

	c.loaded = true
	c.config = config
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

	for _, serviceID := range driverConfig.ServiceIDs {
		for _, service := range c.config.Services {
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
