package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config/consul"
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

	apiVersion, err := c.provisioner.GetValue("usb/api_version")
	if err != nil {
		return nil, err
	}
	config.APIVersion = string(apiVersion)

	config.DriversPath, err = c.GetDriversPath()
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

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
		driverkey = strings.TrimSuffix(driverkey, "/")
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

					driverInstanceInfo, _, err := c.GetDriverInstance(instanceKey)

					if err != nil {
						return nil, err
					}

					dialkeys, err := c.provisioner.GetAllKeys("usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/", "/", nil)

					for _, dialKey := range dialkeys {
						dialKey = strings.TrimSuffix(dialKey, "/")
						dialKey = strings.TrimPrefix(dialKey, "usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/")

						dialInfo, _, err := c.GetDial(dialKey)
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

func (c *consulConfig) GetDriversPath() (string, error) {
	pathValue, _ := c.provisioner.GetValue("usb/drivers_path")

	//TODO fix get value error

	path := string(pathValue)

	if path != "" {
		return path, nil
	}

	if os.Getenv("USB_DRIVER_PATH") != "" {
		path = os.Getenv("USB_DRIVER_PATH")
	} else {
		path = "drivers"
	}

	return path, nil

}

func (c *consulConfig) GetDriver(driverID string) (*Driver, error) {
	var result Driver
	driverType, err := c.provisioner.GetValue("usb/drivers/" + driverID + "/Type")
	if err != nil {
		return &Driver{}, err
	}
	driverName, err := c.provisioner.GetValue("usb/drivers/" + driverID + "/Name")
	if err != nil {
		return &Driver{}, err
	}

	if driverType != nil {
		result.DriverType = string(driverType)
	}
	if driverName != nil {
		result.DriverName = string(driverName)
	}

	instanceKeys, err := c.provisioner.GetAllKeys("usb/drivers/"+driverID+"/instances/", "/", nil)

	for _, instanceKey := range instanceKeys {
		if strings.HasSuffix(instanceKey, "/") {
			instanceKey = strings.TrimSuffix(instanceKey, "/")
			instanceKey = strings.TrimPrefix(instanceKey, "usb/drivers/"+driverID+"/instances/")

			driverInstanceInfo, _, err := c.GetDriverInstance(instanceKey)

			if err != nil {
				return nil, err
			}

			dialkeys, err := c.provisioner.GetAllKeys("usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/", "/", nil)

			for _, dialKey := range dialkeys {
				dialKey = strings.TrimSuffix(dialKey, "/")
				dialKey = strings.TrimPrefix(dialKey, "usb/drivers/"+driverID+"/instances/"+instanceKey+"/dials/")

				dialInfo, _, err := c.GetDial(dialKey)
				if err != nil {
					return nil, err
				}
				if driverInstanceInfo.Dials == nil {
					driverInstanceInfo.Dials = make(map[string]Dial)
				}
				driverInstanceInfo.Dials[dialKey] = *dialInfo
			}
			if result.DriverInstances == nil {
				result.DriverInstances = make(map[string]DriverInstance)
			}
			result.DriverInstances[instanceKey] = *driverInstanceInfo
		}
	}

	return &result, nil
}

func (c *consulConfig) GetDriverInstance(instanceID string) (*DriverInstance, string, error) {
	var instance DriverInstance
	key, err := c.getKey(instanceID)

	if err != nil {
		return nil, "", err
	}
	if key == "" {
		return nil, "", errors.New(fmt.Sprintf("Instance %s not found", instanceID))
	}
	val, err := c.provisioner.GetValue(key + "/Name")
	if err != nil {
		return nil, "", err
	}
	instance.Name = string(val)

	instanceConfig, err := c.provisioner.GetValue(key + "/Configuration")
	if err != nil {
		return nil, "", err
	}

	configuration := json.RawMessage(instanceConfig)

	instance.Configuration = &configuration

	serviceInfo, err := c.provisioner.GetValue(key + "/service")

	err = json.Unmarshal(serviceInfo, &instance.Service)
	if err != nil {
		return nil, "", err
	}
	driverID := strings.Split(key, "/")[2]
	return &instance, driverID, nil
}

func (c *consulConfig) GetService(serviceid string) (*brokerapi.Service, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for _, driver := range config.Drivers {
		for instanceID, instance := range driver.DriverInstances {
			if instance.Service.ID == serviceid {
				return &instance.Service, instanceID, nil
			}
		}
	}
	return nil, "", errors.New(fmt.Sprintf("Service id %s not found", serviceid))
}

func (c *consulConfig) GetDial(dialID string) (*Dial, string, error) {
	var dialInfo Dial
	key, err := c.getDialKey(dialID)
	if err != nil {
		return nil, "", err
	}
	if key == "" {
		return nil, "", errors.New(fmt.Sprintf("Dial key %s not found", dialID))
	}
	data, err := c.provisioner.GetValue(key)
	if err != nil {
		return nil, "", err
	}
	if data == nil {
		return nil, "", errors.New(fmt.Sprintf("Dial %s not found", dialID))
	}

	err = json.Unmarshal(data, &dialInfo)
	instanceID := strings.Split(key, "/")[4]
	return &dialInfo, instanceID, err
}

func (c *consulConfig) SetDriver(driverID string, driver Driver) error {

	err := c.provisioner.AddKV("usb/drivers/"+driverID+"/Type", []byte(driver.DriverType), nil)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV("usb/drivers/"+driverID+"/Name", []byte(driver.DriverName), nil)
	if err != nil {
		return err
	}

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

func (c *consulConfig) getDialKey(dialID string) (string, error) {
	keys, err := c.provisioner.GetAllKeys("usb/drivers/", "", nil)
	if err != nil {
		return "", nil
	}
	for _, key := range keys {
		if strings.Contains(key, dialID) {
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

func (c *consulConfig) DeleteDial(dialID string) error {
	key, err := c.getDialKey(dialID)
	if err != nil {
		return err
	}
	return c.provisioner.DeleteKV(key, nil)
}

func (c *consulConfig) LoadDriverInstance(instanceID string) (*DriverInstance, error) {
	driverInstance, _, err := c.GetDriverInstance(instanceID)
	if err != nil {
		return nil, err
	}

	key, err := c.getKey(instanceID)

	dialkeys, err := c.provisioner.GetAllKeys(key+"/dials/", "/", nil)

	for _, dialKey := range dialkeys {
		dialKey = strings.TrimSuffix(dialKey, "/")
		dialKey = strings.TrimPrefix(dialKey, key+"/dials/")

		dialInfo, _, err := c.GetDial(dialKey)
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

func (c *consulConfig) DriverInstanceNameExists(driverInstanceName string) (bool, error) {
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

			value, err := c.provisioner.GetValue(instance + "Name")
			if err != nil {
				return false, err
			}

			if string(value) == driverInstanceName {
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
		driver = driver + "Type"

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

func (c *consulConfig) DriverExists(driverID string) (bool, error) {
	drivers, err := c.provisioner.GetAllKeys("usb/drivers/", "/", nil)
	if err != nil {
		return false, err
	}
	exists := false
	for _, driver := range drivers {
		if strings.Contains(driver, driverID) {
			exists = true
			break
		}
	}
	return exists, nil
}

func (c *consulConfig) GetPlan(planid string) (*brokerapi.ServicePlan, string, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", "", err
	}
	for _, driver := range config.Drivers {
		for instanceID, instance := range driver.DriverInstances {
			for dialID, dial := range instance.Dials {
				if dial.Plan.ID == planid {
					return &dial.Plan, dialID, instanceID, nil
				}
			}
		}
	}
	return nil, "", "", errors.New(fmt.Sprintf("Plan id %s not found", planid))
}
