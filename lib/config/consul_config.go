package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	"github.com/pivotal-cf/brokerapi"
	"strings"
)

type consulConfig struct {
	address     string
	provisioner consul.ConsulProvisionerInterface
	config      *Config
}

func NewConsulConfig(provisioner consul.ConsulProvisionerInterface) ConfigProvider {
	var consulStruct consulConfig

	consulStruct.provisioner = provisioner

	return &consulStruct
}

func (c *consulConfig) LoadConfiguration() (*Config, error) {
	var config Config

	brokerapiConfig, err := c.provisioner.GetValue("usb/broker_api")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(brokerapiConfig, &config.BrokerAPI)

	if err != nil {
		return nil, err
	}

	managementApiConfig, err := c.provisioner.GetValue("usb/management_api")
	if err != nil {
		return nil, err
	}

	var management ManagementAPI
	err = json.Unmarshal(managementApiConfig, &management)

	config.ManagementAPI = &management

	if err != nil {
		return nil, err
	}

	driverKeys, err := c.provisioner.GetAllKeys("usb/drivers/", "/", nil)
	if err != nil {
		return nil, err
	}
	drivers := make(map[string]Driver)
	for _, driverkey := range driverKeys {
		driverID := strings.TrimPrefix(driverkey, "usb/drivers/")
		if strings.HasSuffix(driverID, "/") == false {

			driverInfo, err := c.GetDriver(driverID)
			if err != nil {
				return nil, err
			}

			instanceKeys, err := c.provisioner.GetAllKeys("usb/drivers/"+driverID+"/instances/", "/", nil)

			for _, instanceKey := range instanceKeys {
				if strings.HasSuffix(instanceKey, "/") {
					instanceKey = strings.TrimSuffix(instanceKey, "/")
					instanceKey = strings.TrimPrefix(instanceKey, "usb/drivers/"+driverID+"/instances/")

					driverInstanceInfo, err := c.GetDriverInstance(instanceKey)

					if err != nil {
						return nil, err
					}

					serv, err := c.GetService(instanceKey)
					if serv != nil {
						driverInstanceInfo.Service = *serv
					}
					if err != nil {
						return nil, err
					}

					dialkeys, err := c.provisioner.GetAllKeys("usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/", "/", nil)

					for _, dialKey := range dialkeys {
						dialKey = strings.TrimSuffix(dialKey, "/")
						dialKey = strings.TrimPrefix(dialKey, "usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/")

						dialInfo, err := c.GetDial(instanceKey, dialKey)
						if err != nil {
							return nil, err
						}
						if driverInstanceInfo.Dials == nil {
							driverInstanceInfo.Dials = make(map[string]Dial)
						}
						driverInstanceInfo.Dials[dialKey] = *dialInfo
					}
					if driverInfo.DriverInstances == nil {
						driverInfo.DriverInstances = make(map[string]DriverInstance)
					}
					driverInfo.DriverInstances[instanceKey] = *driverInstanceInfo
				}
			}

			drivers[driverID] = *driverInfo

		}
	}
	config.Drivers = drivers

	c.config = &config

	return &config, nil
}

func (c *consulConfig) GetDriver(driverID string) (*Driver, error) {
	var result Driver

	val, err := c.provisioner.GetValue("usb/drivers/" + driverID)
	if err != nil {
		return &Driver{}, err
	}
	if val != nil {
		result.DriverType = string(val)
	}
	return &result, nil
}

func (c *consulConfig) GetDriverInstance(instanceID string) (*DriverInstance, error) {
	var instance DriverInstance
	var config json.RawMessage

	key, err := c.getKey(instanceID)
	if err != nil {
		return &DriverInstance{}, err
	}
	if key == "" {
		return &DriverInstance{}, errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}
	val, err := c.provisioner.GetValue(key + "/Name")
	if err != nil {
		return &DriverInstance{}, err
	}
	instance.Name = string(val)

	config, err = c.provisioner.GetValue(key + "/Configuration")
	if err != nil {
		return &DriverInstance{}, err
	}
	instance.Configuration = &config

	return &instance, nil
}

func (c *consulConfig) GetService(instanceID string) (*brokerapi.Service, error) {
	var service brokerapi.Service
	key, err := c.getKey(instanceID)
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}

	val, err := c.provisioner.GetValue(key + "/service")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &service)

	return &service, err
}

func (c *consulConfig) GetDial(instanceID string, dialID string) (*Dial, error) {
	var dialInfo Dial
	key, err := c.getKey(instanceID)
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}

	data, err := c.provisioner.GetValue(key + "/dials/" + dialID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New(fmt.Sprintf("Dial %s not found", dialID))
	}

	err = json.Unmarshal(data, &dialInfo)

	return &dialInfo, err
}

func (c *consulConfig) SetDriver(driverID string, driver Driver) error {

	err := c.provisioner.AddKV("usb/drivers/"+driverID, []byte(driver.DriverType), nil)

	for instanceKey, driverInst := range driver.DriverInstances {
		err = c.SetDriverInstance(driverID, instanceKey, driverInst)
		if err != nil {
			return err
		}
	}

	return err
}

func (c *consulConfig) SetDriverInstance(driverID string, instanceID string, instance DriverInstance) error {

	err := c.provisioner.AddKV("usb/drivers/"+driverID+"/instances/"+instanceID+"/Name", []byte(instance.Name), nil)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV("usb/drivers/"+driverID+"/instances/"+instanceID+"/Configuration", *instance.Configuration, nil)
	if err != nil {
		return err
	}

	err = c.SetService(instanceID, instance.Service)
	if err != nil {
		return err
	}

	for dialKey, dialInfo := range instance.Dials {
		err = c.SetDial(instanceID, dialKey, dialInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *consulConfig) SetService(instanceID string, service brokerapi.Service) error {

	key, err := c.getKey(instanceID)

	if err != nil {
		return err
	}
	if key == "" {
		return errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}

	data, err := json.Marshal(service)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV(key+"/service", data, nil)

	return err
}

func (c *consulConfig) SetDial(instanceID string, dialID string, dial Dial) error {
	key, err := c.getKey(instanceID)
	if err != nil {
		return err
	}
	if key == "" {
		return errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}

	data, err := json.Marshal(dial)
	if err != nil {
		return err
	}
	err = c.provisioner.AddKV(key+"/dials/"+dialID, data, nil)

	return err
}

func (c *consulConfig) getKey(instanceID string) (string, error) {
	keys, err := c.provisioner.GetAllKeys("usb/drivers/", "", nil)
	if err != nil {
		return "", nil
	}
	for _, key := range keys {
		if strings.Contains(key, instanceID) {
			key = strings.TrimSuffix(key, "/Configuration")
			key = strings.TrimSuffix(key, "/Name")
			return key, nil
		}
	}
	return "", nil
}

func (c *consulConfig) DeleteDriver(driverID string) error {
	return c.provisioner.DeleteKVs("usb/drivers/"+driverID, nil)
}

func (c *consulConfig) DeleteDriverInstance(instanceID string) error {
	key, err := c.getKey(instanceID)
	if err != nil {
		return err
	}
	return c.provisioner.DeleteKVs(key, nil)
}

func (c *consulConfig) DeleteService(instanceID string) error {
	key, err := c.getKey(instanceID)
	if err != nil {
		return err
	}
	return c.provisioner.DeleteKV(key+"/service", nil)
}

func (c *consulConfig) DeleteDial(instanceID string, dialID string) error {
	key, err := c.getKey(instanceID)
	if err != nil {
		return err
	}
	return c.provisioner.DeleteKV(key+"/dials/"+dialID, nil)
}

func (c *consulConfig) LoadDriverInstance(instanceID string) (*DriverInstance, error) {

	driverInstance, err := c.GetDriverInstance(instanceID)
	service, err := c.GetService(instanceID)

	driverInstance.Service = *service

	key, err := c.getKey(instanceID)

	dialkeys, err := c.provisioner.GetAllKeys(key+"/dials/", "/", nil)

	for _, dialKey := range dialkeys {
		dialKey = strings.TrimSuffix(dialKey, "/")
		dialKey = strings.TrimPrefix(dialKey, key+"/dials/")

		dialInfo, err := c.GetDial(instanceID, dialKey)
		if err != nil {
			return nil, err
		}
		if driverInstance.Dials == nil {
			driverInstance.Dials = make(map[string]Dial)
		}
		driverInstance.Dials[dialKey] = *dialInfo
	}

	return driverInstance, err
}

func (c *consulConfig) GetUaaAuthConfig() (*UaaAuth, error) {
	conf := (*json.RawMessage)(c.config.ManagementAPI.Authentication)

	uaa := Uaa{}
	err := json.Unmarshal(*conf, &uaa)
	if err != nil {
		return nil, err
	}
	return &uaa.UaaAuth, nil
}

func (c *consulConfig) ServiceNameExists(serviceName string) (bool, error) {
	drivers, err := c.provisioner.GetAllKeys("usb/drivers/", "/", nil)
	if err != nil {
		return false, err
	}

	for _, driver := range drivers {

		instances, err := c.provisioner.GetAllKeys(driver+"/instances/", "/", nil)
		if err != nil {
			return false, err
		}

		for _, instance := range instances {

			var service brokerapi.Service

			value, err := c.provisioner.GetValue(instance + "service")
			if err != nil {
				return false, err
			}

			err = json.Unmarshal(value, &service)
			if err != nil {
				return false, err
			}

			if service.Name == serviceName {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *consulConfig) DriverTypeExists(driverType string) (bool, error) {
	drivers, err := c.provisioner.GetAllKeys("usb/drivers/", "/", nil)
	if err != nil {
		return false, err
	}

	for _, driver := range drivers {
		driver = strings.TrimSuffix(driver, "/")

		value, err := c.provisioner.GetValue(driver)
		if err != nil {
			return false, err
		}

		if string(value) == driverType {
			return true, nil
		}
	}

	return false, nil
}
