package config

import (
	"encoding/json"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config/redis"
)

const usbKey = "usb"

type redisConfig struct {
	provider redis.Provisioner
}

//NewRedisConfig generates and returns a new redis config provider
func NewRedisConfig(provider redis.Provisioner) Provider {
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

func (c *redisConfig) LoadDriverInstance(driverInstanceID string) (*Instance, error) {
	driver, _, err := c.GetInstance(driverInstanceID)
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

func (c *redisConfig) SetInstance(instanceID string, instance Instance) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		instanceInfo := &instance
		if _, ok := configuration.Instances[instanceID]; ok {
			configuration.Instances[instanceID] = *instanceInfo
		} else {
			if configuration.Instances == nil {
				configuration.Instances = make(map[string]Instance)
			}
			configuration.Instances[instanceID] = *instanceInfo
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

func (c *redisConfig) GetInstance(instanceID string) (*Instance, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for diKey, i := range config.Instances {
		if diKey == instanceID {
			return &i, diKey, nil
		}
	}
	return nil, "", nil
}

func (c *redisConfig) DeleteInstance(instanceID string) error {

	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	if _, ok := config.Instances[instanceID]; ok {
		delete(config.Instances, instanceID)
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
		if _, ok := configuration.Instances[instanceID]; ok {
			instance := configuration.Instances[instanceID]
			instance.Service = service
			configuration.Instances[instanceID] = instance
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

	for instanceID, instance := range config.Instances {
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

	if instance, ok := config.Instances[instanceID]; ok {
		instance.Service = brokerapi.Service{}
		config.Instances[instanceID] = instance
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
		if instance, ok := configuration.Instances[instanceID]; ok {
			if _, ok := instance.Dials[dialID]; ok {
				instance.Dials[dialID] = *dialDetails
			} else {
				if instance.Dials == nil {
					instance.Dials = make(map[string]Dial)
				}
				instance.Dials[dialID] = *dialDetails
				configuration.Instances[instanceID] = instance
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

	for instanceID, instance := range config.Instances {
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

	for instanceID, instance := range config.Instances {
		if _, ok := instance.Dials[dialID]; ok {
			delete(instance.Dials, dialID)
			config.Instances[instanceID] = instance
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

func (c *redisConfig) InstanceNameExists(driverInstanceName string) (bool, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return false, err
	}

	for _, di := range config.Instances {
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
	for iID, i := range config.Instances {
		for dialID, di := range i.Dials {
			if di.Plan.ID == planid {
				return &di.Plan, dialID, iID, nil
			}
		}
	}

	return nil, "", "", nil
}
