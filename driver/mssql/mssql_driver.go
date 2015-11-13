package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mssql/config"
	"github.com/hpcloud/cf-usb/driver/mssql/driverdata"
	"github.com/hpcloud/cf-usb/driver/mssql/mssqlprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

const happySqlPasswordPolicySuffix = "Aa_0"

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

type MssqlDriver struct {
	logger lager.Logger
	conf   config.MssqlDriverConfig
	db     mssqlprovisioner.MssqlProvisionerInterface
}

func NewMssqlDriver(logger lager.Logger, db mssqlprovisioner.MssqlProvisionerInterface) driver.Driver {
	return &MssqlDriver{
		logger: logger,
		db:     db,
	}
}

func (e *MssqlDriver) init(conf *json.RawMessage) error {
	mssqlConfig := config.MssqlDriverConfig{}
	err := json.Unmarshal(*conf, &mssqlConfig)
	e.logger.Info("Mssql Driver initializing")

	var mssqlConConfig = map[string]string{}
	mssqlConConfig["server"] = mssqlConfig.Host
	mssqlConConfig["port"] = strconv.Itoa(mssqlConfig.Port)
	mssqlConConfig["user id"] = mssqlConfig.User
	mssqlConConfig["password"] = mssqlConfig.Pass

	err = e.db.Connect("mssql", mssqlConConfig)
	if err != nil {
		return err
	}

	e.conf = mssqlConfig
	return nil
}

func (e *MssqlDriver) Ping(request *json.RawMessage, response *bool) error {
	*response = false

	err := e.init(request)
	if err != nil {
		return err
	}

	*response = true
	return nil
}

func (e *MssqlDriver) GetDailsSchema(empty string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (e *MssqlDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (e *MssqlDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID

	err = e.db.CreateDatabase(databaseName)
	if err != nil {
		return err
	}

	response.Status = status.Created
	return nil
}

func (e *MssqlDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	response.Status = status.DoesNotExist
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID

	created, err := e.db.IsDatabaseCreated(databaseName)
	if err != nil {
		return err
	}

	if created {
		response.Status = status.Exists
	}

	return nil
}

func (e *MssqlDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID
	username := request.CredentialsID

	randomString, err := secureRandomString(32)
	if err != nil {
		return err
	}

	password := randomString + happySqlPasswordPolicySuffix

	err = e.db.CreateUser(databaseName, username, password)
	if err != nil {
		return err
	}
	data := MssqlBindingCredentials{
		Host:     e.conf.Host,
		Port:     e.conf.Port,
		Username: username,
		Password: password,
		ConnectionString: generateConnectionString(e.conf.Host, e.conf.Port,
			databaseName, username, password),
	}

	*response = data
	return nil
}

func (e *MssqlDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.DoesNotExist
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID
	username := request.CredentialsID

	created, err := e.db.IsUserCreated(databaseName, username)
	if err != nil {
		return err
	}

	if created {
		response.Status = status.Exists
	}

	return nil
}

func (e *MssqlDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID
	username := request.CredentialsID

	err = e.db.DeleteUser(databaseName, username)
	if err != nil {
		return err
	}

	response.Status = status.Deleted
	return nil
}

func (e *MssqlDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	databaseName := e.conf.DbIdentifierPrefix + request.InstanceID

	err = e.db.DeleteDatabase(databaseName)
	if err != nil {
		return err
	}

	response.Status = status.Deleted

	return nil
}
