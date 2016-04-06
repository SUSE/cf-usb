package redis

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

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
	return &RedisDriver{logger: logger.Session("redis-driver"), redisProvisioner: provisioner}
}

func (d *RedisDriver) init(conf *json.RawMessage) error {
	d.logger.Info("init-driver", lager.Data{"configValue": string(*conf)})

	redisConfig := config.RedisDriverConfig{}

	err := json.Unmarshal(*conf, &redisConfig)
	if err != nil {
		return err
	}

	err = d.redisProvisioner.Connect(redisConfig)
	if err != nil {
		return err
	}

	d.conf = redisConfig

	return nil
}

func (d *RedisDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

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
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	response.Description = "Error creating instance"

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.redisProvisioner.CreateContainer(request.InstanceID)
	if err != nil {
		return err
	}

	response.Status = status.Created
	response.InstanceID = request.InstanceID
	response.Description = "Instance created"

	return nil
}

func (d *RedisDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	response.Description = "Error deleting instance"

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.redisProvisioner.DeleteContainer(request.InstanceID)
	if err != nil {
		return err
	}

	response.Status = status.Deleted
	response.InstanceID = request.InstanceID
	response.Description = "Instance deleted"

	return nil
}

func (d *RedisDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	cred, err := d.redisProvisioner.GetCredentials(request.InstanceID)
	if err != nil {
		return err
	}

	host := ""
	dockerUrl, err := url.Parse(d.conf.DockerEndpoint)
	if err != nil {
		return err
	}

	if dockerUrl.Scheme == "unix" {
		host, err = getLocalIP()
		if err != nil {
			return err
		}
	} else {
		host = strings.Split(dockerUrl.Host, ":")[0]
	}
	data := RedisBindingCredentials{
		Password: cred["password"],
		Port:     cred["port"],
		Host:     host,
		Hostname: host,
		Uri:      fmt.Sprintf("redis://:%s@%s:%s/", cred["password"], host, cred["port"]),
	}

	*response = data

	return nil
}

func (d *RedisDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	response.Status = status.Deleted

	return nil
}

func (d *RedisDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	response.Status = status.DoesNotExist

	return nil
}

func (d *RedisDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *RedisDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *RedisDriver) GetParametersSchema(request string, response *string) error {
	//Does not support custom parameters
	return nil
}

func (d *RedisDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

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
