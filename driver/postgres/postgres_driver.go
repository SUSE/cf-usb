package postgres

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-golang/lager"
)

type postgresDriver struct {
	driverProperties    model.DriverInitRequest
	defaultConnParams   postgresprovisioner.PostgresServiceProperties
	logger              lager.Logger
	postgresProvisioner postgresprovisioner.PostgresProvisionerInterface
}

func NewPostgresDriver(logger lager.Logger) driver.Driver {
	return &postgresDriver{logger: logger}
}

func (driver *postgresDriver) Init(driverProperties model.DriverInitRequest, response *string) error {
	driver.driverProperties = driverProperties

	conf := (*json.RawMessage)(driverProperties.DriverConfig)
	dsp := postgresprovisioner.PostgresServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	driver.logger.Info("init-driver", lager.Data{"user": dsp.User, "password": dsp.Password, "host": dsp.Host, "port": dsp.Port, "dbname": dsp.Dbname, "sslmode": dsp.Sslmode})

	driver.defaultConnParams = dsp
	driver.postgresProvisioner = postgresprovisioner.NewPostgresProvisioner(dsp, driver.logger)
	driver.postgresProvisioner.Init()
	if err != nil {
		driver.logger.Fatal("error-initializing-provisioner", err)
		return err
	}

	*response = "Sucessfully initialized postgres driver"
	return nil
}

func (driver *postgresDriver) Ping(request string, response *bool) error {
	return nil
}

func (driver *postgresDriver) GetDailsSchema(request string, response *string) error {
	dailsSchema, err := data.Asset("schemas/dails.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (driver *postgresDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := data.Asset("scehmas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (driver *postgresDriver) ProvisionInstance(request model.ProvisionInstanceRequest, response *bool) error {
	err := driver.postgresProvisioner.CreateDatabase(request.InstanceID)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
		return err
	}

	*response = true
	return nil
}

func (driver *postgresDriver) InstanceExists(instanceID string, response *bool) error {
	exist, err := driver.postgresProvisioner.DatabaseExists(instanceID)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
	}
	response = &exist
	return nil
}

func (driver *postgresDriver) GenerateCredentials(request model.CredentialsRequest, response *interface{}) error {

	username := request.InstanceID + request.CredentialsID
	password, err := secureRandomString(32)
	if err != nil {
		return err
	}

	err = driver.postgresProvisioner.CreateUser(request.InstanceID, username, password)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
		return err
	}

	data := PostgresBindingCredentials{
		Hostname:         driver.defaultConnParams.Host,
		Name:             request.InstanceID,
		Password:         password,
		Port:             driver.defaultConnParams.Port,
		Username:         username,
		ConnectionString: generateConnectionString(driver.defaultConnParams.Host, driver.defaultConnParams.Port, request.InstanceID, username, password),
	}
	*response = data
	return nil
}

func (driver *postgresDriver) CredentialsExist(request model.CredentialsRequest, response *bool) error {
	username := request.InstanceID + request.CredentialsID

	exist, err := driver.postgresProvisioner.UserExists(username)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
	}
	response = &exist
	return nil
}

func (driver *postgresDriver) RevokeCredentials(request model.CredentialsRequest, response *interface{}) error {
	driver.logger.Info("unbind-request", lager.Data{"credentialsID": request.CredentialsID, "InstanceID": request.InstanceID})
	username := request.InstanceID + request.CredentialsID

	err := driver.postgresProvisioner.DeleteUser(request.InstanceID, username)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
		return err
	}
	*response = ""

	return nil
}
func (driver *postgresDriver) DeprovisionInstance(instanceID string, response *interface{}) error {
	driver.logger.Info("deprovision-request", lager.Data{"instance-id": instanceID})

	err := driver.postgresProvisioner.DeleteDatabase(instanceID)
	if err != nil {
		driver.logger.Fatal("provision-error", err)
		return err
	}

	*response = "Successfully deprovisoned"
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
