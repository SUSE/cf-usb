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
	for _, d := range c.config.Drivers {
		if d.ID == driverID {
			return d, nil
		}
	}
	return Driver{}, errors.New(fmt.Sprintf("Driver ID: %s not found", driverID))
}

func (c *fileConfig) GetDriverInstance(instanceID string) (DriverInstance, error) {
	for _, d := range c.config.Drivers {
		for _, i := range d.DriverInstances {
			if i.ID == instanceID {
				return *i, nil
			}
		}
	}
	return DriverInstance{}, errors.New(fmt.Sprintf("Driver Instance ID: %s not found", instanceID))
}

func (c *fileConfig) GetService(instanceID string) (brokerapi.Service, error) {
	instance, err := c.GetDriverInstance(instanceID)
	if err != nil {
		return brokerapi.Service{}, errors.New(fmt.Sprintf("Driver Instance ID: %s not found", instanceID))
	}

	return instance.Service, nil
}

func (c *fileConfig) GetDial(instanceID string, dialID string) (Dial, error) {
	instance, err := c.GetDriverInstance(instanceID)
	if err != nil {
		return Dial{}, errors.New(fmt.Sprintf("Driver Instance ID: %s not found", instanceID))
	}
	for _, dial := range instance.Dials {
		if dial.ID == dialID {
			return dial, nil
		}
	}
	return Dial{}, errors.New(fmt.Sprintf("Dial ID: %s not found for Driver Instance ID:%s", dialID, instanceID))
}

func (c *fileConfig) SetDriver(driverInfo Driver) error {
	return errors.New("SetDriver not available for file config provider")
}

func (c *fileConfig) SetDriverInstance(driverID string, instance DriverInstance) error {
	return errors.New("SetDriverInstance not available for file config provider")
}

func (c *fileConfig) SetService(instanceID string, service brokerapi.Service) error {
	return errors.New("SetService not available for file config provider")
}

func (c *fileConfig) SetDial(instanceID string, dialInfo Dial) error {
	return errors.New("SetDial not available for file config provider")
}

func (c *fileConfig) DeleteDriver(driverID string) error {
	return errors.New("DeleteDriver not available for file config provider")
}

func (c *fileConfig) DeleteDriverInstance(instanceID string) error {
	return errors.New("DeleteDriverInstance not available for file config provider")
}

func (c *fileConfig) DeleteService(instanceID string) error {
	return errors.New("DeleteService not available for file config provider")
}

func (c *fileConfig) DeleteDial(instanceID string, dialID string) error {
	return errors.New("DeleteDial not available for file config provider")
}

func parseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
