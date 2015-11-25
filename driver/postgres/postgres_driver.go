package postgres

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/postgres/config"
	"github.com/hpcloud/cf-usb/driver/postgres/driverdata"
	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type PostgresDriver struct {
	conf                config.PostgresDriverConfig
	logger              lager.Logger
	postgresProvisioner postgresprovisioner.PostgresProvisionerInterface
}

func NewPostgresDriver(logger lager.Logger, provisioner postgresprovisioner.PostgresProvisionerInterface) driver.Driver {
	return &PostgresDriver{logger: logger, postgresProvisioner: provisioner}
}

func (d *PostgresDriver) init(conf *json.RawMessage) error {

	postgressConfig := config.PostgresDriverConfig{}
	err := json.Unmarshal(*conf, &postgressConfig)
	d.logger.Info("Postgress Driver initializing")
	err = d.postgresProvisioner.Connect(postgressConfig)
	if err != nil {
		return err
	}
	d.conf = postgressConfig
	return nil
}

func (d *PostgresDriver) Ping(request *json.RawMessage, response *bool) error {
	*response = false
	err := d.init(request)
	if err != nil {
		return err
	}
	*response = true
	return nil
}

func (d *PostgresDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *PostgresDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *PostgresDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	err = d.postgresProvisioner.CreateDatabase(request.InstanceID)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}

	response.Status = status.Created
	return nil
}

func (d *PostgresDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist
	exist, err := d.postgresProvisioner.DatabaseExists(request.InstanceID)
	if err != nil {
		d.logger.Fatal("get-instance-error", err)
		return err
	}
	if exist {
		response.Status = status.Exists
	}

	return nil
}

func (d *PostgresDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	username := request.InstanceID + request.CredentialsID
	password, err := secureRandomString(32)
	if err != nil {
		return err
	}

	err = d.postgresProvisioner.CreateUser(request.InstanceID, username, password)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}

	data := PostgresBindingCredentials{
		Hostname:         d.conf.Host,
		Name:             request.InstanceID,
		Password:         password,
		Port:             d.conf.Port,
		Username:         username,
		ConnectionString: generateConnectionString(d.conf.Host, d.conf.Port, request.InstanceID, username, password),
	}
	*response = data
	return nil
}

func (d *PostgresDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist

	username := request.InstanceID + request.CredentialsID

	exist, err := d.postgresProvisioner.UserExists(username)
	if err != nil {
		d.logger.Fatal("provision-error", err)
	}
	if exist {
		response.Status = status.Exists
	}
	return nil
}

func (d *PostgresDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}

	d.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})
	username := request.InstanceID + request.CredentialsID

	err = d.postgresProvisioner.DeleteUser(request.InstanceID, username)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}
	response.Status = status.Deleted

	return nil
}
func (d *PostgresDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := d.init(request.Config)
	if err != nil {
		return err
	}
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request.InstanceID})

	err = d.postgresProvisioner.DeleteDatabase(request.InstanceID)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}

	response.Status = status.Deleted
	return nil
}

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}
