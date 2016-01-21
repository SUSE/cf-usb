package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

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

	if configuration.DriversPath == "" {
		if os.Getenv("USB_DRIVER_PATH") != "" {
			configuration.DriversPath = os.Getenv("USB_DRIVER_PATH")
		} else {
			configuration.DriversPath = "drivers"
		}
	}

	return &configuration, nil
}

func (c *redisConfig) GetDriversPath() (string, error) {
	config, err := c.LoadConfiguration()

	if err != nil {
		return "", err
	}

	return config.DriversPath, nil
}

func (c *redisConfig) LoadDriverInstance(driverInstanceID string) (*DriverInstance, error) {
	driver, err := c.GetDriverInstance(driverInstanceID)
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
func (c *redisConfig) SetDriver(driverID string, driver Driver) error {

	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	updated := false
	for dId, _ := range config.Drivers {
		if dId == driverID {
			config.Drivers[dId] = driver
			updated = true
			break
		}
	}
	if config.Drivers == nil {
		config.Drivers = make(map[string]Driver)
	}
	if !updated {
		config.Drivers[driverID] = driver
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

func (c *redisConfig) GetDriver(driverID string) (*Driver, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for dID, d := range config.Drivers {
		if dID == driverID {
			return &d, nil
		}
	}

	return nil, nil
}

func (c *redisConfig) DeleteDriver(driverID string) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	for dID, _ := range config.Drivers {
		if dID == driverID {
			delete(config.Drivers, driverID)
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

func (c *redisConfig) SetDriverInstance(driverID string, instanceID string, instance DriverInstance) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		instanceInfo := &instance
		if driverInfo, ok := configuration.Drivers[driverID]; ok {
			if _, ok := driverInfo.DriverInstances[instanceID]; ok {
				driverInfo.DriverInstances[instanceID] = *instanceInfo
				configuration.Drivers[driverID] = driverInfo
			} else {
				if driverInfo.DriverInstances == nil {
					driverInfo.DriverInstances = make(map[string]DriverInstance)
				}
				driverInfo.DriverInstances[instanceID] = *instanceInfo
				configuration.Drivers[driverID] = driverInfo
			}
			data, err := json.Marshal(configuration)
			if err != nil {
				return err
			}

			err = c.provider.SetKV(usbKey, string(data), 0)
			if err != nil {
				return err
			}

		} else {
			return errors.New(fmt.Sprintf("Driver id %s not found", driverID))
		}
	}

	return nil
}

func (c *redisConfig) GetDriverInstance(instanceID string) (*DriverInstance, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, d := range config.Drivers {
		for diKey, i := range d.DriverInstances {
			if diKey == instanceID {
				return &i, nil
			}
		}
	}
	return nil, nil
}

func (c *redisConfig) DeleteDriverInstance(instanceID string) error {

	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	for _, d := range config.Drivers {
		delete(d.DriverInstances, instanceID)
		break
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
		for driverKey, driverInfo := range configuration.Drivers {
			if _, ok := driverInfo.DriverInstances[instanceID]; ok {
				instance := driverInfo.DriverInstances[instanceID]
				instance.Service = service
				configuration.Drivers[driverKey].DriverInstances[instanceID] = instance
				break
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

func (c *redisConfig) GetService(serviceID string) (*brokerapi.Service, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for _, d := range config.Drivers {
		for instanceID, instance := range d.DriverInstances {
			if instance.Service.ID == serviceID {
				return &instance.Service, instanceID, nil
			}
		}
	}
	return nil, "", errors.New(fmt.Sprintf("Service id %s not found", serviceID))
}

func (c *redisConfig) DeleteService(instanceID string) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	for _, d := range config.Drivers {
		if instance, ok := d.DriverInstances[instanceID]; ok {
			instance.Service = brokerapi.Service{}
			d.DriverInstances[instanceID] = instance
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

func (c *redisConfig) SetDial(instanceID string, dialID string, dial Dial) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}
	if config != nil {
		configuration := *config
		dialDetails := &dial
		for driverKey, driverInfo := range configuration.Drivers {
			if instance, ok := driverInfo.DriverInstances[instanceID]; ok {
				if _, ok := instance.Dials[dialID]; ok {
					instance.Dials[dialID] = *dialDetails
				} else {
					if instance.Dials == nil {
						instance.Dials = make(map[string]Dial)
					}
					instance.Dials[dialID] = *dialDetails
					configuration.Drivers[driverKey].DriverInstances[instanceID] = instance
				}
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

func (c *redisConfig) GetDial(dialID string) (*Dial, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, d := range config.Drivers {
		for _, instance := range d.DriverInstances {
			if dialInfo, ok := instance.Dials[dialID]; ok {
				return &dialInfo, nil
			}
		}
	}

	return nil, nil
}

func (c *redisConfig) DeleteDial(dialID string) error {
	config, err := c.LoadConfiguration()
	if err != nil {
		return err
	}

	for _, d := range config.Drivers {
		for instanceID, instance := range d.DriverInstances {
			if _, ok := instance.Dials[dialID]; ok {
				delete(instance.Dials, dialID)
				d.DriverInstances[instanceID] = instance
				break
			}
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

func (c *redisConfig) ServiceNameExists(serviceName string) (bool, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return false, err
	}

	for _, d := range config.Drivers {
		for _, di := range d.DriverInstances {
			if di.Service.Name == serviceName {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *redisConfig) DriverTypeExists(driverType string) (bool, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return false, err
	}

	for _, d := range config.Drivers {
		if d.DriverType == driverType {
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
	for _, d := range config.Drivers {
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
