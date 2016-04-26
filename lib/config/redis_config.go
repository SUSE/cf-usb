package config

import (
	"encoding/json"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config/redis"
)

const usbKey = "usb"

type redisConfig struct {
	provider redis.RedisProvisionerInterface
}

func NewRedisConfig(provider redis.RedisProvisionerInterface) ConfigProvider {
	provisioner := redisConfig{}
	provisioner.provider = provider
	return &provisioner
}

func (c *redisConfig) LoadConfiguration() (*Config, error) {

	var configuration Config

	configurationValue, err := c.provider.GetValue(usbKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(configurationValue), &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (c *redisConfig) LoadDriverInstance(driverInstanceID string) (*DriverInstance, error) {
	driver, _, err := c.GetDriverInstance(driverInstanceID)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (c *redisConfig) GetUaaAuthConfig() (*UaaAuth, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	conf := (*json.RawMessage)(config.ManagementAPI.Authentication)

	uaa := Uaa{}
	err = json.Unmarshal(*conf, &uaa)
	if err != nil {
		return nil, err
	}
	return &uaa.UaaAuth, nil
}

func (c *redisConfig) SetDriverInstance(instanceID string, instance DriverInstance) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		instanceInfo := &instance
		if _, ok := configuration.DriverInstances[instanceID]; ok {
			configuration.DriverInstances[instanceID] = *instanceInfo
		} else {
			if configuration.DriverInstances == nil {
				configuration.DriverInstances = make(map[string]DriverInstance)
			}
			configuration.DriverInstances[instanceID] = *instanceInfo
		}
		data, err := json.Marshal(configuration)
		if err != nil {
			return err
		}

		err = c.provider.SetKV(usbKey, string(data), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisConfig) GetDriverInstance(instanceID string) (*DriverInstance, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for diKey, i := range config.DriverInstances {
		if diKey == instanceID {
			return &i, diKey, nil
		}
	}
	return nil, "", nil
}

func (c *redisConfig) DeleteDriverInstance(instanceID string) error {

	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	if _, ok := config.DriverInstances[instanceID]; ok {
		delete(config.DriverInstances, instanceID)
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = c.provider.SetKV(usbKey, string(data), 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisConfig) SetService(instanceID string, service brokerapi.Service) error {

	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		if _, ok := configuration.DriverInstances[instanceID]; ok {
			instance := configuration.DriverInstances[instanceID]
			instance.Service = service
			configuration.DriverInstances[instanceID] = instance
		}

		data, err := json.Marshal(configuration)
		if err != nil {
			return err
		}

		err = c.provider.SetKV(usbKey, string(data), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisConfig) GetService(serviceID string) (*brokerapi.Service, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for instanceID, instance := range config.DriverInstances {
		if instance.Service.ID == serviceID {
			return &instance.Service, instanceID, nil
		}
	}

	return nil, "", nil
}

func (c *redisConfig) DeleteService(instanceID string) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	if instance, ok := config.DriverInstances[instanceID]; ok {
		instance.Service = brokerapi.Service{}
		config.DriverInstances[instanceID] = instance
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = c.provider.SetKV(usbKey, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) SetDial(instanceID string, dialID string, dial Dial) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		dialDetails := &dial
		if instance, ok := configuration.DriverInstances[instanceID]; ok {
			if _, ok := instance.Dials[dialID]; ok {
				instance.Dials[dialID] = *dialDetails
			} else {
				if instance.Dials == nil {
					instance.Dials = make(map[string]Dial)
				}
				instance.Dials[dialID] = *dialDetails
				configuration.DriverInstances[instanceID] = instance
			}
		}

		data, err := json.Marshal(configuration)

		if err != nil {
			return err
		}

		err = c.provider.SetKV(usbKey, string(data), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *redisConfig) GetDial(dialID string) (*Dial, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for instanceID, instance := range config.DriverInstances {
		if dialInfo, ok := instance.Dials[dialID]; ok {
			return &dialInfo, instanceID, nil
		}
	}

	return nil, "", nil
}

func (c *redisConfig) DeleteDial(dialID string) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	for instanceID, instance := range config.DriverInstances {
		if _, ok := instance.Dials[dialID]; ok {
			delete(instance.Dials, dialID)
			config.DriverInstances[instanceID] = instance
			break
		}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = c.provider.SetKV(usbKey, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) DriverInstanceNameExists(driverInstanceName string) (bool, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return false, err
	}

	for _, di := range config.DriverInstances {
		if di.Name == driverInstanceName {
			return true, nil
		}
	}

	return false, nil
}

func (c *redisConfig) GetPlan(planid string) (*brokerapi.ServicePlan, string, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", "", err
	}
	for iID, i := range config.DriverInstances {
		for dialID, di := range i.Dials {
			if di.Plan.ID == planid {
				return &di.Plan, dialID, iID, nil
			}
		}
	}

	return nil, "", "", nil
}
