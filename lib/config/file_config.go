package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SUSE/cf-usb/lib/brokermodel"
)

type fileConfig struct {
	path   string
	loaded bool
	config *Config
}

//NewFileConfig builds and returns a new file config Provider
func NewFileConfig(path string) Provider {
	return &fileConfig{path: path, loaded: false}
}

func (c *fileConfig) SaveConfiguration(config Config, overwrite bool) error {
	return fmt.Errorf("Not implemented")
}

func (c *fileConfig) LoadConfiguration() (*Config, error) {
	var config *Config

	jsonConf, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}

	config, err = parseJSON(jsonConf)
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

func (c *fileConfig) LoadDriverInstance(instanceID string) (*Instance, error) {

	var instance Instance
	exists := false

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, err
		}
	}

	for diKey, di := range c.config.Instances {
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

func (c *fileConfig) GetInstance(instanceID string) (*Instance, string, error) {
	for diKey, i := range c.config.Instances {
		if diKey == instanceID {
			return &i, diKey, nil
		}
	}
	return nil, "", nil
}

func (c *fileConfig) GetService(serviceID string) (*brokermodel.CatalogService, string, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, "", err
		}
	}
	for diKey, i := range c.config.Instances {
		if i.Service.ID == serviceID {
			return &i.Service, diKey, nil
		}
	}

	return nil, "", nil
}

func (c *fileConfig) GetDial(dialID string) (*Dial, string, error) {
	for instanceID, instance := range c.config.Instances {
		for dialID, dial := range instance.Dials {
			if dialID == dialID {
				return &dial, instanceID, nil
			}
		}

	}
	return nil, "", nil
}

func (c *fileConfig) SetInstance(instanceID string, instance Instance) error {
	return errors.New("SetDriverInstance not available for file config provider")
}

func (c *fileConfig) SetService(instanceID string, service brokermodel.CatalogService) error {
	return errors.New("SetService not available for file config provider")
}

func (c *fileConfig) SetDial(instanceID string, dialID string, dialInfo Dial) error {
	return errors.New("SetDial not available for file config provider")
}

func (c *fileConfig) DeleteInstance(instanceID string) error {
	return errors.New("DeleteDriverInstance not available for file config provider")
}

func (c *fileConfig) DeleteService(instanceID string) error {
	return errors.New("DeleteService not available for file config provider")
}

func (c *fileConfig) DeleteDial(dialID string) error {
	return errors.New("DeleteDial not available for file config provider")
}

func (c *fileConfig) InstanceNameExists(driverInstanceName string) (bool, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return false, err
		}
	}

	for _, di := range c.config.Instances {
		if di.Name == driverInstanceName {
			return true, nil
		}
	}

	return false, nil
}

func (c *fileConfig) GetPlan(planid string) (*brokermodel.Plan, string, string, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, "", "", err
		}
	}
	for iID, i := range c.config.Instances {
		for dialID, di := range i.Dials {
			if di.Plan.ID == planid {
				return &di.Plan, dialID, iID, nil
			}
		}

	}
	return nil, "", "", nil
}

func parseJSON(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
