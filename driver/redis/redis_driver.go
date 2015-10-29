package redis

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/redis/config"
	"github.com/hpcloud/cf-usb/driver/redis/driverdata"
	"github.com/hpcloud/cf-usb/driver/redis/redisprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"

	"github.com/pivotal-golang/lager"
)

type RedisDriver struct {
	conf             config.RedisDriverConfig
	redisProvisioner redisprovisioner.RedisProvisionerInterface
	logger           lager.Logger
}

func NewRedisDriver(logger lager.Logger, provisioner redisprovisioner.RedisProvisionerInterface) driver.Driver {
	return &RedisDriver{logger: logger, redisProvisioner: provisioner}
}

func (d *RedisDriver) init(conf *json.RawMessage) error {

	redisConfig := config.RedisDriverConfig{}
	err := json.Unmarshal(*conf, &redisConfig)
	d.logger.Info("Postgress Driver initializing")
	err = d.redisProvisioner.Connect(redisConfig)
	if err != nil {
		return err
	}
	d.conf = redisConfig
	return nil
}

func (d *RedisDriver) Ping(request *json.RawMessage, response *bool) error {
	err := d.init(request)
	if err != nil {
		return err
	}

	err = d.redisProvisioner.PingServer()
	if err != nil {
		*response = false
		return err
	}
	*response = true
	return nil
}

func (d *RedisDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}
	err = d.redisProvisioner.CreateContainer(request.InstanceID)
	if err != nil {
		return err
	}
	response.Status = status.Created
	return nil
}

func (d *RedisDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.redisProvisioner.DeleteContainer(request.InstanceID)
	if err != nil {
		return err
	}

	response.Status = status.Deleted
	return nil
}

func (d *RedisDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	cred, err := d.redisProvisioner.GetCredentials(request.InstanceID)
	if err != nil {
		return err
	}

	localIp, err := getLocalIP()
	if err != nil {
		return err
	}

	data := RedisBindingCredentials{
		Password: cred["password"],
		Port:     cred["port"],
		Host:     localIp,
	}

	*response = data
	return nil
}

func (d *RedisDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.Deleted
	return nil
}

func (d *RedisDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.DoesNotExist
	return nil
}

func (d *RedisDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *RedisDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *RedisDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}
	response.Status = status.DoesNotExist
	exists, err := d.redisProvisioner.ContainerExists(request.InstanceID)
	if err != nil {
		return err
	}

	if exists {
		response.Status = status.Exists
	}
	return nil
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, inface := range interfaces {
		addresses, err := inface.Addrs()
		if err != nil {
			return "", err
		}
		for _, address := range addresses {
			ipnet, ok := address.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP.To4()
			if ip == nil || ip.IsLoopback() {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("Could not find IP address")
}
