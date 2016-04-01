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
	return &RabbitmqDriver{logger: logger.Session("rabbitmq-driver"), rabbitmqProvisioner: provisioner}
}

func (d *RabbitmqDriver) init(conf *json.RawMessage) error {
	d.logger.Info("init-driver", lager.Data{"configValue": string(*conf)})

	rabbitmqConfig := config.RabbitmqDriverConfig{}

	err := json.Unmarshal(*conf, &rabbitmqConfig)
	if err != nil {
		d.logger.Error("init-driver-failed", err)
		return err
	}

	err = d.rabbitmqProvisioner.Connect(rabbitmqConfig)
	if err != nil {
		d.logger.Error("init-driver-failed", err)
		return err
	}

	d.conf = rabbitmqConfig

	return nil
}

func (d *RabbitmqDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

	err := d.init(request)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.PingServer()
	if err != nil {
		*response = false
		d.logger.Error("ping-request-failed", err)
		return err
	}

	*response = true

	return nil
}

func (d *RabbitmqDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *RabbitmqDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *RabbitmqDriver) GetParametersSchema(request string, response *string) error {
	//Does not support custom parameters
	return nil
}

func (d *RabbitmqDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	response.Description = "Error creating instance"

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.CreateContainer(request.InstanceID)
	if err != nil {
		d.logger.Error("provision-instance-request-failed", err)
		return err
	}

	response.Status = status.Created
	response.InstanceID = request.InstanceID
	response.Description = "Instance created"

	return nil
}

func (d *RabbitmqDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist

	exists, err := d.rabbitmqProvisioner.ContainerExists(request.InstanceID)
	if err != nil {
		d.logger.Error("get-instance-request-failed", err)
		return err
	}

	if exists {
		response.Status = status.Exists
	}

	return nil
}

func (d *RabbitmqDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	credentials, err := d.rabbitmqProvisioner.CreateUser(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("generate-credentials-request-failed", err)
		return err
	}

	data := RabbitmqBindingCredentials{
		Hostname:     credentials["host"],
		Host:         credentials["host"],
		VHost:        credentials["vhost"],
		Port:         credentials["port"],
		Username:     credentials["user"],
		Password:     credentials["password"],
		Uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", credentials["user"], credentials["password"], credentials["host"], credentials["port"], credentials["vhost"]),
		DashboardUrl: fmt.Sprintf("http://%s:%s", credentials["host"], credentials["mgmt_port"]),
	}

	*response = data

	return nil
}

func (d *RabbitmqDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist

	exist, err := d.rabbitmqProvisioner.UserExists(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("credentials-exists-request-failed", err)
		return err
	}
	if exist {
		response.Status = status.Exists
	}

	return nil
}

func (d *RabbitmqDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.DeleteUser(request.InstanceID, request.CredentialsID)
	if err != nil {
		d.logger.Error("revoke-credentials-request-failed", err)
		return err
	}

	response.Status = status.Deleted

	return nil
}
func (d *RabbitmqDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	response.Description = "Error deleting instance"

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.rabbitmqProvisioner.DeleteContainer(request.InstanceID)
	if err != nil {
		d.logger.Error("deprovision-request-failed", err)
		return err
	}

	response.Status = status.Deleted
	response.InstanceID = request.InstanceID
	response.Description = "Instance deleted"

	return nil
}
