package postgres

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"regexp"

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
	return &PostgresDriver{logger: logger.Session("postgres-driver"), postgresProvisioner: provisioner}
}

func (d *PostgresDriver) init(conf *json.RawMessage) error {
	d.logger.Info("init-driver", lager.Data{"configValue": string(*conf)})

	postgressConfig := config.PostgresDriverConfig{}

	err := json.Unmarshal(*conf, &postgressConfig)

	err = d.postgresProvisioner.Connect(postgressConfig)
	if err != nil {
		return err
	}

	d.conf = postgressConfig

	return nil
}

func (d *PostgresDriver) Ping(request *json.RawMessage, response *bool) error {
	d.logger.Info("ping-request", lager.Data{"request": string(*request)})

	*response = false

	err := d.init(request)
	if err != nil {
		return err
	}

	*response = true

	return nil
}

func (d *PostgresDriver) GetDailsSchema(request string, response *string) error {
	d.logger.Info("get-dails-schema-request", lager.Data{"request": request})

	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (d *PostgresDriver) GetConfigSchema(request string, response *string) error {
	d.logger.Info("get-config-schema-request", lager.Data{"request": request})

	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (d *PostgresDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("provision-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config), "dials": string(*request.Dials)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	dbName, err := d.getMD5Hash(request.InstanceID)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}

	err = d.postgresProvisioner.CreateDatabase(dbName)
	if err != nil {
		d.logger.Fatal("provision-error", err)
		return err
	}

	response.Status = status.Created

	return nil
}

func (d *PostgresDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	d.logger.Info("get-instance-request", lager.Data{"instance-id": request.InstanceID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	dbName, err := d.getMD5Hash(request.InstanceID)
	if err != nil {
		d.logger.Fatal("get-instance-request-failed", err)
		return err
	}

	response.Status = status.DoesNotExist
	exist, err := d.postgresProvisioner.DatabaseExists(dbName)
	if err != nil {
		d.logger.Fatal("get-instance-request-failed", err)
		return err
	}
	if exist {
		response.Status = status.Exists
	}

	return nil
}

func (d *PostgresDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	d.logger.Info("generate-credentials-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	dbName, err := d.getMD5Hash(request.InstanceID)
	if err != nil {
		d.logger.Fatal("generate-credentials-request-failed", err)
		return err
	}

	username, err := d.getMD5Hash(request.InstanceID + request.CredentialsID)
	if err != nil {
		d.logger.Fatal("generate-credentials-request-failed", err)
		return err
	}

	password, err := secureRandomString(32)
	if err != nil {
		return err
	}

	err = d.postgresProvisioner.CreateUser(dbName, username, password)
	if err != nil {
		d.logger.Fatal("generate-credentials-request-failed", err)
		return err
	}

	data := PostgresBindingCredentials{
		Hostname:         d.conf.Host,
		Host:             d.conf.Host,
		Database:         dbName,
		Password:         password,
		Port:             d.conf.Port,
		Username:         username,
		ConnectionString: generateConnectionString(d.conf.Host, d.conf.Port, dbName, username, password),
	}

	*response = data

	return nil
}

func (d *PostgresDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("credentials-exists-request", lager.Data{"instance-id": request.InstanceID, "credentials-id": request.CredentialsID, "config": string(*request.Config)})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	response.Status = status.DoesNotExist

	username, err := d.getMD5Hash(request.InstanceID + request.CredentialsID)
	if err != nil {
		d.logger.Fatal("credentials-exists-request-failed", err)
		return err
	}

	exist, err := d.postgresProvisioner.UserExists(username)
	if err != nil {
		d.logger.Fatal("credentials-exists-request-failed", err)
	}
	if exist {
		response.Status = status.Exists
	}

	return nil
}

func (d *PostgresDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	d.logger.Info("revoke-credentials-request", lager.Data{"credentials-id": request.CredentialsID, "instance-id": request.InstanceID})

	err := d.init(request.Config)
	if err != nil {
		return err
	}
	dbName, err := d.getMD5Hash(request.InstanceID)
	if err != nil {
		d.logger.Fatal("revoke-credentials-request-failed", err)
		return err
	}

	username, err := d.getMD5Hash(request.InstanceID + request.CredentialsID)
	if err != nil {
		d.logger.Fatal("revoke-credentials-request-failed", err)
		return err
	}

	err = d.postgresProvisioner.DeleteUser(dbName, username)
	if err != nil {
		d.logger.Fatal("revoke-credentials-request-failed", err)
		return err
	}

	response.Status = status.Deleted

	return nil
}
func (d *PostgresDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	d.logger.Info("deprovision-request", lager.Data{"instance-id": request})

	err := d.init(request.Config)
	if err != nil {
		return err
	}

	dbName, err := d.getMD5Hash(request.InstanceID)
	if err != nil {
		d.logger.Fatal("deprovision-request-failed", err)
		return err
	}

	err = d.postgresProvisioner.DeleteDatabase(dbName)
	if err != nil {
		d.logger.Fatal("deprovision-request-failed", err)
		return err
	}

	response.Status = status.Deleted

	return nil
}

func (d *PostgresDriver) getMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}
