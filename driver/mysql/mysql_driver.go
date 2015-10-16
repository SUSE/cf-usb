package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mysql/driverdata"
	"github.com/hpcloud/cf-usb/driver/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-golang/lager"
)

type MysqlDriver struct {
	User   string `json:"userid"`
	Pass   string `json:"password"`
	Host   string `json:"server"`
	Port   string `json:"port"`
	db     mysqlprovisioner.MysqlProvisionerInterface
	logger lager.Logger
}

func (e *MysqlDriver) secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		e.logger.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func NewMysqlDriver(logger lager.Logger) driver.Driver {
	return &MysqlDriver{logger: logger}
}

func (e *MysqlDriver) Init(configuration model.DriverInitRequest, response *string) error {
	err := json.Unmarshal(*configuration.DriverConfig, &e)
	e.logger.Info("Mysql Driver initializing")
	e.db, err = mysqlprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port, e.logger)
	return err
}

func (e *MysqlDriver) Ping(empty string, result *bool) error {

	return nil
}

func (e *MysqlDriver) GetDailsSchema(empty string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (e *MysqlDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}
func (e *MysqlDriver) ProvisionInstance(request model.ProvisionInstanceRequest, result *bool) error {
	err := e.db.CreateDatabase(strings.Replace(request.InstanceID, "-", "", -1))
	if err != nil {
		return err
	}
	*result = true
	return nil
}

func (e *MysqlDriver) InstanceExists(instanceID string, result *bool) error {
	created, err := e.db.IsDatabaseCreated(strings.Replace(instanceID, "-", "", -1))
	if err != nil {
		return err
	}
	*result = created
	return nil
}

func (e *MysqlDriver) GenerateCredentials(request model.CredentialsRequest, response *interface{}) error {
	username := strings.Replace(request.CredentialsID, "-", "", -1)
	if len(username) > 16 {
		username = username[:16]
	}
	password := e.secureRandomString(32)

	err := e.db.CreateUser(strings.Replace(request.InstanceID, "-", "", -1), username, password)
	if err != nil {
		return err
	}
	data := MysqlBindingCredentials{
		Host:             e.Host,
		Port:             e.Port,
		Username:         username,
		Password:         password,
		ConnectionString: generateConnectionString(e.Host, e.Port, strings.Replace(request.InstanceID, "-", "", -1), username, password),
	}

	*response = data
	return nil
}

func (e *MysqlDriver) CredentialsExist(request model.CredentialsRequest, response *bool) error {
	username := strings.Replace(request.CredentialsID, "-", "", -1)
	if len(username) > 16 {
		username = username[:16]
	}
	created, err := e.db.IsUserCreated(username)
	if err != nil {
		return err
	}
	*response = created
	return nil
}

func (e *MysqlDriver) RevokeCredentials(request model.CredentialsRequest, response *interface{}) error {
	username := strings.Replace(request.CredentialsID, "-", "", -1)
	if len(username) > 16 {
		username = username[:16]
	}
	err := e.db.DeleteUser(username)
	if err != nil {
		return err
	}
	*response = fmt.Sprintf("Credentials for %s revoked", username)
	return nil
}

func (e *MysqlDriver) DeprovisionInstance(instanceID string, response *interface{}) error {
	err := e.db.DeleteDatabase(strings.Replace(instanceID, "-", "", -1))
	if err != nil {
		return err
	}
	*response = fmt.Sprintf("Deprovisioned %s", instanceID)

	return nil
}
