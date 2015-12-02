package config

import (
	"encoding/json"
	"github.com/hpcloud/cf-usb/lib/config/redis"
	"github.com/pivotal-cf/brokerapi"
)

type redisConfig struct {
	provider redis.RedisProvisionerInterface
	config   *Config
}

func NewRedisConfig(provider redis.RedisProvisionerInterface) ConfigProvider {
	provisioner := redisConfig{}
	provisioner.provider = provider
	return &provisioner
}

func (c *redisConfig) LoadConfiguration() (*Config, error) {
	var configuration Config

	value, err := c.provider.GetValue("broker_api")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &configuration.BrokerAPI)

	if err != nil {
		return nil, err
	}

	value, err = c.provider.GetValue("management_api")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &configuration.ManagementAPI)

	if err != nil {
		return nil, err
	}

	value, err = c.provider.GetValue("log_level")

	if err != nil {
		return nil, err
	}

	configuration.LogLevel = value
	exists, err := c.provider.KeyExists("drivers")
	if err != nil {
		return nil, err
	}
	if exists {
		value, err = c.provider.GetValue("drivers")

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(value), &configuration.Drivers)
		if err != nil {
			return nil, err
		}
	}
	return &configuration, nil
}

func (c *redisConfig) LoadDriverInstance(driverInstanceID string) (*DriverInstance, error) {

	return nil, nil
}

func (c *redisConfig) GetUaaAuthConfig() (*UaaAuth, error) {
	conf := (*json.RawMessage)(c.config.ManagementAPI.Authentication)

	uaa := Uaa{}
	err := json.Unmarshal(*conf, &uaa)
	if err != nil {
		return nil, err
	}
	return &uaa.UaaAuth, nil
}
func (c *redisConfig) SetDriver(driver Driver) error {

	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		if d.ID == driver.ID {
			d = driver
			break
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) GetDriver(driverID string) (Driver, error) {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return Driver{}, err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		if d.ID == driverID {
			return d, nil
		}
	}

	return Driver{}, nil
}

func (c *redisConfig) DeleteDriver(driverID string) error {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for i, d := range c.config.Drivers {
		if d.ID == driverID {
			c.config.Drivers = append(c.config.Drivers[:i], c.config.Drivers[i+1:]...)
			break
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) SetDriverInstance(driverID string, instance DriverInstance) error {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}
	modified := false
	for _, d := range c.config.Drivers {
		if d.ID == driverID {
			for _, instanceInfo := range d.DriverInstances {
				if instanceInfo.ID == instance.ID {
					instanceInfo = &instance
					modified = true
					break
				}
			}
			if !modified {
				d.DriverInstances = append(d.DriverInstances, &instance)
			}
			break
		}
	}

	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisConfig) GetDriverInstance(instanceID string) (DriverInstance, error) {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return DriverInstance{}, err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, i := range d.DriverInstances {
			if i.ID == instanceID {
				return *i, nil
			}
		}
	}
	return DriverInstance{}, nil
}

func (c *redisConfig) DeleteDriverInstance(instanceID string) error {

	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for i, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				d.DriverInstances = append(d.DriverInstances[:i], d.DriverInstances[i+1:]...)
				break
			}
		}
	}

	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisConfig) SetService(instanceID string, service brokerapi.Service) error {

	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				instanceInfo.Service = service
				break
			}
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *redisConfig) GetService(instanceID string) (brokerapi.Service, error) {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return brokerapi.Service{}, err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				return instanceInfo.Service, nil
			}
		}
	}
	return brokerapi.Service{}, nil
}

func (c *redisConfig) DeleteService(instanceID string) error {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				instanceInfo.Service = brokerapi.Service{}
				break
			}
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) SetDial(instanceID string, dial Dial) error {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}
	modified := false
	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				for _, dialInfo := range instanceInfo.Dials {
					if dialInfo.ID == dial.ID {
						dialInfo = dial
						modified = true
						break
					}
				}
				if !modified {
					instanceInfo.Dials = append(instanceInfo.Dials, dial)
					break
				}
			}
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *redisConfig) GetDial(instanceID string, dialID string) (Dial, error) {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return Dial{}, err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				for _, dialInfo := range instanceInfo.Dials {
					if dialInfo.ID == dialID {
						return dialInfo, nil
					}
				}
			}
		}
	}

	return Dial{}, nil
}

func (c *redisConfig) DeleteDial(instanceID string, dialID string) error {
	if c.config == nil {
		config, err := c.LoadConfiguration()
		if err != nil {
			return err
		}
		c.config = config
	}

	for _, d := range c.config.Drivers {
		for _, instanceInfo := range d.DriverInstances {
			if instanceInfo.ID == instanceID {
				for i, dialInfo := range instanceInfo.Dials {
					if dialInfo.ID == dialID {
						instanceInfo.Dials = append(instanceInfo.Dials[:i], instanceInfo.Dials[i+1:]...)
						break
					}
				}
				break
			}
		}
	}
	data, err := json.Marshal(c.config.Drivers)
	if err != nil {
		return err
	}

	err = c.provider.SetKV("drivers", string(data), 0)
	if err != nil {
		return err
	}
	return nil
}
