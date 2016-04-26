package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/frodenas/brokerapi"
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

func (c *fileConfig) LoadDriverInstance(instanceID string) (*DriverInstance, error) {

	var instance DriverInstance
	exists := false

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, err
		}
	}

	for diKey, di := range c.config.DriverInstances {
		if diKey == instanceID {
			instance = di
			exists = true
			break
		}
	}

	if !exists {
		return nil, nil
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

func (c *fileConfig) GetDriverInstance(instanceID string) (*DriverInstance, string, error) {
	for diKey, i := range c.config.DriverInstances {
		if diKey == instanceID {
			return &i, diKey, nil
		}
	}
	return nil, "", nil
}

func (c *fileConfig) GetService(serviceID string) (*brokerapi.Service, string, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, "", err
		}
	}
	for diKey, i := range c.config.DriverInstances {
		if i.Service.ID == serviceID {
			return &i.Service, diKey, nil
		}
	}

	return nil, "", nil
}

func (c *fileConfig) GetDial(dialID string) (*Dial, string, error) {
	for instanceID, instance := range c.config.DriverInstances {
		for dialID, dial := range instance.Dials {
			if dialID == dialID {
				return &dial, instanceID, nil
			}
		}

	}
	return nil, "", nil
}

func (c *fileConfig) SetDriverInstance(instanceID string, instance DriverInstance) error {
	return errors.New("SetDriverInstance not available for file config provider")
}

func (c *fileConfig) SetService(instanceID string, service brokerapi.Service) error {
	return errors.New("SetService not available for file config provider")
}

func (c *fileConfig) SetDial(instanceID string, dialID string, dialInfo Dial) error {
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

func (c *fileConfig) DeleteDial(dialID string) error {
	return errors.New("DeleteDial not available for file config provider")
}

func (c *fileConfig) DriverInstanceNameExists(driverInstanceName string) (bool, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return false, err
		}
	}

	for _, di := range c.config.DriverInstances {
		if di.Name == driverInstanceName {
			return true, nil
		}
	}

	return false, nil
}

func (c *fileConfig) GetPlan(planid string) (*brokerapi.ServicePlan, string, string, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, "", "", err
		}
	}
	for iID, i := range c.config.DriverInstances {
		for dialID, di := range i.Dials {
			if di.Plan.ID == planid {
				return &di.Plan, dialID, iID, nil
			}
		}

	}
	return nil, "", "", nil
}

func parseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
