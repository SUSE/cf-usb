package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/config"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/driverdata"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/rabbitmqprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type RabbitmqDriver struct {
	conf                config.RabbitmqDriverConfig
	logger              lager.Logger
	rabbitmqProvisioner rabbitmqprovisioner.RabbitmqProvisionerInterface
}

func NewRabbitmqDriver(logger lager.Logger, provisioner rabbitmqprovisioner.RabbitmqProvisionerInterface) driver.Driver {
	return &RabbitmqDriver{logger: logger, rabbitmqProvisioner: provisioner}
}

func (d *RabbitmqDriver) init(conf *json.RawMessage) error {

	rabbitmqConfig := config.RabbitmqDriverConfig{}
	err := json.Unmarshal(*conf, &rabbitmqConfig)
	d.logger.Info("Rabbitmq Driver initializing")
	err = d.rabbitmqProvisioner.Connect(rabbitmqConfig)
	if err != nil {
		d.logger.Error("Error initializing RabbitMQ driver", err)
		return err
	}
	d.conf = rabbitmqConfig
	return nil
}

func (d *RabbitmqDriver) Ping(request *json.RawMessage, response *bool) error {
	err := d.init(request)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.PingServer()
	if err != nil {
		*response = false
		d.logger.Error("ping-error", err)
		return err
	}
	*response = true
	return nil
}

func (d *RabbitmqDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *RabbitmqDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *RabbitmqDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}
	err = d.rabbitmqProvisioner.CreateContainer(request.InstanceID)
	if err != nil {
		d.logger.Error("provision-error", err)
		return err
	}
	response.Status = status.Created
	return nil
}

func (d *RabbitmqDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}
	response.Status = status.DoesNotExist
	exists, err := d.rabbitmqProvisioner.ContainerExists(request.InstanceID)
	if err != nil {
		d.logger.Error("get-instance-error", err)
		return err
	}

	if exists {
		response.Status = status.Exists
	}
	return nil
}

func (d *RabbitmqDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	credentials, err := d.rabbitmqProvisioner.CreateUser(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("genetate-credentials-error", err)
		return err
	}

	data := RabbitmqBindingCredentials{
		Host:         credentials["host"],
		VHost:        credentials["vhost"],
		Port:         credentials["port"],
		Username:     credentials["user"],
		Password:     credentials["password"],
		Uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", credentials["user"], credentials["password"], credentials["host"], credentials["port"], credentials["vhost"]),
		DashboardUrl: fmt.Sprintf("http://%s:%s", credentials["host"], credentials["port"]),
	}
	*response = data
	return nil
}

func (d *RabbitmqDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist

	exist, err := d.rabbitmqProvisioner.UserExists(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("get-credentials-error", err)
	}
	if exist {
		response.Status = status.Exists
	}
	return nil
}

func (d *RabbitmqDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	d.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})

	err = d.rabbitmqProvisioner.DeleteUser(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("revoke-credentials-error", err)
		return err
	}
	response.Status = status.Deleted

	return nil
}
func (d *RabbitmqDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.DeleteContainer(request.InstanceID)
	if err != nil {
		d.logger.Error("unprovision-error", err)
		return err
	}

	response.Status = status.Deleted
	return nil
}
