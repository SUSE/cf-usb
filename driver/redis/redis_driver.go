package redis

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/redis/driverdata"
	"github.com/hpcloud/cf-usb/driver/redis/redisprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"

	"github.com/pivotal-golang/lager"
)

type redisDriver struct {
	driverProperties  model.DriverInitRequest
	defaultConnParams redisprovisioner.RedisServiceProperties
	redisProvisioner  redisprovisioner.RedisProvisionerInterface
	logger            lager.Logger
}

func NewRedisDriver(logger lager.Logger) driver.Driver {
	return &redisDriver{logger: logger}
}

func (driver *redisDriver) Init(driverProperties model.DriverInitRequest, response *string) error {
	driver.driverProperties = driverProperties

	conf := (*json.RawMessage)(driverProperties.DriverConfig)

	serviceProperties := redisprovisioner.RedisServiceProperties{}

	err := json.Unmarshal(*conf, &serviceProperties)

	driver.defaultConnParams = serviceProperties
	driver.redisProvisioner = redisprovisioner.NewRedisProvisioner(serviceProperties, driver.logger)
	driver.redisProvisioner.Init()
	if err != nil {
		driver.logger.Fatal("error-initializing-provisioner", err)
		return err
	}

	*response = "Sucessfully initialized redis driver"
	return nil
}

func (driver *redisDriver) ProvisionInstance(request model.ProvisionInstanceRequest, response *bool) error {
	err := driver.redisProvisioner.CreateContainer(request.InstanceID)
	if err != nil {
		*response = false
		return err
	}
	*response = true
	return nil
}

func (driver *redisDriver) DeprovisionInstance(request string, response *interface{}) error {
	err := driver.redisProvisioner.DeleteContainer(request)
	if err != nil {
		*response = false
		return err
	}
	*response = true
	return nil
}

func (driver *redisDriver) GenerateCredentials(request model.CredentialsRequest, response *interface{}) error {

	cred, err := driver.redisProvisioner.GetCredentials(request.InstanceID)
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

func (driver *redisDriver) RevokeCredentials(request model.CredentialsRequest, response *interface{}) error {
	return nil
}

func (driver *redisDriver) CredentialsExist(request model.CredentialsRequest, response *bool) error {
	exists, err := driver.redisProvisioner.ContainerExists(request.InstanceID)
	if err != nil {
		return err
	}
	*response = exists
	return nil
}

func (driver *redisDriver) Ping(request string, response *bool) error {
	err := driver.redisProvisioner.PingServer()
	if err != nil {
		*response = false
		return err
	}
	*response = true
	return nil
}

func (driver *redisDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (driver *redisDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (driver *redisDriver) InstanceExists(request string, response *bool) error {
	exists, err := driver.redisProvisioner.ContainerExists(request)
	if err != nil {
		return err
	}
	*response = exists
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
