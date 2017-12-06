package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/SUSE/cf-usb/lib/brokermodel"
	"github.com/SUSE/cf-usb/lib/config/consul"
)

type consulConfig struct {
	address     string
	provisioner consul.Provisioner
	config      *Config
}

//NewConsulConfig builds and returns a new consul ConfigProvider
func NewConsulConfig(provisioner consul.Provisioner) Provider {
	var consulStruct consulConfig

	consulStruct.provisioner = provisioner

	return &consulStruct
}

func (c *consulConfig) SaveConfiguration(config Config, overwrite bool) error {
	return fmt.Errorf("Not implemented")
}

func (c *consulConfig) LoadConfiguration() (*Config, error) {
	var config Config

	apiVersion, err := c.provisioner.GetValue("usb/api_version")
	if err != nil {
		return nil, err
	}
	config.APIVersion = string(apiVersion)

	brokerapiConfig, err := c.provisioner.GetValue("usb/broker_api")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(brokerapiConfig, &config.BrokerAPI)

	if err != nil {
		return nil, err
	}

	managementAPIConfig, err := c.provisioner.GetValue("usb/management_api")
	if err != nil {
		return nil, err
	}

	var management ManagementAPI
	err = json.Unmarshal(managementAPIConfig, &management)
	if err != nil {
		return nil, err
	}

	config.ManagementAPI = &management

	if err != nil {
		return nil, err
	}

	instanceKeys, err := c.provisioner.GetAllKeys("usb/instances/", "/", nil)
	if err != nil {
		return nil, err
	}

	instances := make(map[string]Instance)
	for _, instanceKey := range instanceKeys {
		instanceKey = strings.TrimSuffix(instanceKey, "/")
		instanceKey := strings.TrimPrefix(instanceKey, "usb/instances/")

		if strings.HasSuffix(instanceKey, "/") == false {

			driverInstanceInfo, _, err := c.GetInstance(instanceKey)

			if err != nil {
				return nil, err
			}

			dialkeys, err := c.provisioner.GetAllKeys("usb/instances/"+instanceKey+"/dials/", "/", nil)

			for _, dialKey := range dialkeys {
				dialKey = strings.TrimSuffix(dialKey, "/")
				dialKey = strings.TrimPrefix(dialKey, "usb/instances/"+instanceKey+"/dials/")

				dialInfo, _, err := c.GetDial(dialKey)
				if err != nil {
					return nil, err
				}
				if driverInstanceInfo.Dials == nil {
					driverInstanceInfo.Dials = make(map[string]Dial)
				}
				driverInstanceInfo.Dials[dialKey] = *dialInfo
			}

			instances[instanceKey] = *driverInstanceInfo

		}
	}
	config.Instances = instances

	c.config = &config

	return &config, nil
}

func (c *consulConfig) GetInstance(instanceID string) (*Instance, string, error) {
	var instance Instance
	key, err := c.getKey(instanceID)

	if err != nil {
		return nil, "", err
	}
	if key == "" {
		return nil, "", nil
	}
	val, err := c.provisioner.GetValue(key + "/name")
	if err != nil {
		return nil, "", err
	}

	instance.Name = string(val)

	target, err := c.provisioner.GetValue(key + "/target_url")
	if err != nil {
		return nil, "", err
	}
	instance.TargetURL = string(target)

	authKey, err := c.provisioner.GetValue(key + "/authentication_key")
	if err != nil {
		return nil, "", err
	}
	instance.AuthenticationKey = string(authKey)

	serviceInfo, err := c.provisioner.GetValue(key + "/service")

	err = json.Unmarshal(serviceInfo, &instance.Service)
	if err != nil {
		return nil, "", err
	}

	caCert, err := c.provisioner.GetValue(key + "/ca_cert")
	if err != nil {
		return nil, "", err
	}
	instance.CaCert = string(caCert)

	skipSsl, err := c.provisioner.GetValue(key + "/skip_ssl")
	if err != nil {
		return nil, "", err
	}

	skipSslBool, err := strconv.ParseBool(string(skipSsl))
	if err != nil {
		return nil, "", err
	}

	instance.SkipSsl = skipSslBool

	driverID := strings.Split(key, "/")[2]
	return &instance, driverID, nil
}

func (c *consulConfig) GetService(serviceid string) (*brokermodel.CatalogService, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", err
	}

	for instanceID, instance := range config.Instances {
		if instance.Service.ID == serviceid {
			return &instance.Service, instanceID, nil
		}
	}
	return nil, "", nil
}

func (c *consulConfig) GetDial(dialID string) (*Dial, string, error) {
	var dialInfo Dial
	key, err := c.getDialKey(dialID)
	if err != nil {
		return nil, "", err
	}
	if key == "" {
		return nil, "", nil
	}
	data, err := c.provisioner.GetValue(key)
	if err != nil {
		return nil, "", err
	}
	if data == nil {
		return nil, "", nil
	}

	err = json.Unmarshal(data, &dialInfo)
	instanceID := strings.Split(key, "/")[4]
	return &dialInfo, instanceID, err
}

func (c *consulConfig) SetInstance(instanceID string, instance Instance) error {

	err := c.provisioner.AddKV("usb/instances/"+instanceID+"/name", []byte(instance.Name), nil)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV("usb/instances/"+instanceID+"/target_url", []byte(instance.TargetURL), nil)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV("usb/instances/"+instanceID+"/authentication_key", []byte(instance.AuthenticationKey), nil)
	if err != nil {
		return err
	}

	err = c.provisioner.AddKV("usb/instances/"+instanceID+"/ca_cert", []byte(instance.CaCert), nil)
	if err != nil {
		return err
	}

	skipSsl := strconv.FormatBool(instance.SkipSsl)
	err = c.provisioner.AddKV("usb/instances/"+instanceID+"/skip_ssl", []byte(skipSsl), nil)
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

func (c *consulConfig) SetService(instanceID string, service brokermodel.CatalogService) error {

	key, err := c.getKey(instanceID)

	if err != nil {
		return err
	}
	if key == "" {
		return fmt.Errorf("Instance %s not found", instanceID)
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
		return fmt.Errorf("Instance %s not found", instanceID)
	}

	data, err := json.Marshal(dial)
	if err != nil {
		return err
	}
	err = c.provisioner.AddKV(key+"/dials/"+dialID, data, nil)

	return err
}

func (c *consulConfig) getKey(instanceID string) (string, error) {
	keys, err := c.provisioner.GetAllKeys("usb/instances/", "", nil)
	if err != nil {
		return "", nil
	}

	for _, key := range keys {
		if strings.Contains(key, instanceID) {
			key = strings.TrimSuffix(key, "/name")
			key = strings.TrimSuffix(key, "/authentication_key")
			key = strings.TrimSuffix(key, "/target_url")
			key = strings.TrimSuffix(key, "/service")
			return key, nil
		}
	}
	return "", nil
}

func (c *consulConfig) getDialKey(dialID string) (string, error) {
	keys, err := c.provisioner.GetAllKeys("usb/instances/", "", nil)
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

func (c *consulConfig) DeleteInstance(instanceID string) error {
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

func (c *consulConfig) LoadDriverInstance(instanceID string) (*Instance, error) {
	driverInstance, _, err := c.GetInstance(instanceID)
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

func (c *consulConfig) InstanceNameExists(driverInstanceName string) (bool, error) {
	instances, err := c.provisioner.GetAllKeys("usb/instances/", "/", nil)
	if err != nil {
		return false, err
	}
	for _, instance := range instances {
		value, err := c.provisioner.GetValue(instance + "name")
		if err != nil {
			return false, err
		}
		if strings.ToLower(string(value)) == strings.ToLower(driverInstanceName) {
			return true, nil
		}
	}

	return false, nil
}

func (c *consulConfig) GetPlan(planid string) (*brokermodel.Plan, string, string, error) {
	config, err := c.LoadConfiguration()
	if err != nil {
		return nil, "", "", err
	}
	for instanceID, instance := range config.Instances {
		for dialID, dial := range instance.Dials {
			if dial.Plan.ID == planid {
				return &dial.Plan, dialID, instanceID, nil
			}
		}
	}
	return nil, "", "", nil
}
