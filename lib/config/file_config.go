package config

import (
	"encoding/json"
	"errors"
	"fmt"
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

	c.config.DriversPath, err = c.GetDriversPath()
	if err != nil {
		return nil, err
	}

	if os.Getenv("PORT") != "" {
		c.config.BrokerAPI.Listen = ":" + os.Getenv("PORT")
	}
	return config, nil
}

func (c *fileConfig) GetDriversPath() (string, error) {

	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return "", err
		}
	}

	path := c.config.DriversPath

	if path == "" {
		if os.Getenv("USB_DRIVER_PATH") != "" {
			path = os.Getenv("USB_DRIVER_PATH")
		} else {
			path = "drivers"
		}
	}

	return path, nil

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

	for _, d := range c.config.Drivers {
		for diKey, di := range d.DriverInstances {
			if diKey == instanceID {
				instance = di
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

func (c *fileConfig) GetDriver(driverID string) (*Driver, error) {
	for id, d := range c.config.Drivers {
		if id == driverID {
			return &d, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Driver ID: %s not found", driverID))
}

func (c *fileConfig) GetDriverInstance(instanceID string) (*DriverInstance, error) {
	for _, d := range c.config.Drivers {
		for diKey, i := range d.DriverInstances {
			if diKey == instanceID {
				return &i, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Driver Instance ID: %s not found", instanceID))
}

func (c *fileConfig) GetService(serviceID string) (*brokerapi.Service, string, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return nil, "", err
		}
	}
	for _, d := range c.config.Drivers {
		for diKey, i := range d.DriverInstances {
			if i.Service.ID == serviceID {
				return &i.Service, diKey, nil
			}
		}
	}

	return nil, "", errors.New(fmt.Sprintf("Service id %s not found", serviceID))
}

func (c *fileConfig) GetDial(dialID string) (*Dial, error) {
	for _, d := range c.config.Drivers {
		for _, instance := range d.DriverInstances {
			for dialID, dial := range instance.Dials {
				if dialID == dialID {
					return &dial, nil
				}
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Dial ID: %s not found", dialID))
}

func (c *fileConfig) SetDriver(driverID string, driverInfo Driver) error {
	return errors.New("SetDriver not available for file config provider")
}

func (c *fileConfig) SetDriverInstance(driverID string, instanceID string, instance DriverInstance) error {
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

func (c *fileConfig) ServiceNameExists(serviceName string) (bool, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return false, err
		}
	}

	for _, d := range c.config.Drivers {
		for _, di := range d.DriverInstances {
			if di.Service.Name == serviceName {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *fileConfig) DriverTypeExists(driverType string) (bool, error) {
	if !c.loaded {
		_, err := c.LoadConfiguration()
		if err != nil {
			return false, err
		}
	}

	for _, d := range c.config.Drivers {
		if d.DriverType == driverType {
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
	for _, d := range c.config.Drivers {
		for iID, i := range d.DriverInstances {
			for dialID, di := range i.Dials {
				if di.Plan.ID == planid {
					return &di.Plan, dialID, iID, nil
				}
			}
		}
	}
	return nil, "", "", errors.New(fmt.Sprintf("Plan id %s not found", planid))
}

func parseJson(jsonConf []byte) (*Config, error) {
	config := &Config{}

	err := json.Unmarshal(jsonConf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
