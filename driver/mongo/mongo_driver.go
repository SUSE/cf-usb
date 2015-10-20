package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mongo/driverdata"
	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-golang/lager"
)

type MongoDriver struct {
	User   string `json:"user id"`
	Pass   string `json:"password"`
	Host   string `json:"server"`
	Port   string `json:"port"`
	db     mongoprovisioner.MongoProvisionerInterface
	logger lager.Logger
}

func (e *MongoDriver) secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		e.logger.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func NewMongoDriver(logger lager.Logger) driver.Driver {
	return &MongoDriver{logger: logger}
}

func (e *MongoDriver) Init(configuration model.DriverInitRequest, response *string) error {
	err := json.Unmarshal(*configuration.DriverConfig, &e)
	e.logger.Info("Mongo Driver initializing")
	e.db, err = mongoprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port, e.logger)
	return err
}

func (e *MongoDriver) Ping(empty string, result *bool) error {
	_, err := mongoprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port, e.logger)
	if err != nil {
		*result = false
		return err
	}
	*result = true
	return nil
}

func (e *MongoDriver) GetDailsSchema(empty string, response *string) error {
	dailsSchema, err := driverdata.Asset("schemas/dials.json")
	if err != nil {
		return err
	}

	*response = string(dailsSchema)

	return nil
}

func (e *MongoDriver) GetConfigSchema(request string, response *string) error {
	configSchema, err := driverdata.Asset("schemas/config.json")
	if err != nil {
		return err
	}

	*response = string(configSchema)

	return nil
}

func (e *MongoDriver) ProvisionInstance(request model.ProvisionInstanceRequest, result *bool) error {
	err := e.db.CreateDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	return nil
}

func (e *MongoDriver) InstanceExists(instanceID string, result *bool) error {
	created, err := e.db.IsDatabaseCreated(instanceID)
	if err != nil {
		return err
	}
	*result = created
	return nil
}

func (e *MongoDriver) GenerateCredentials(request model.CredentialsRequest, response *interface{}) error {
	username := request.InstanceID + "-" + request.CredentialsID
	password := e.secureRandomString(32)

	err := e.db.CreateUser(request.InstanceID, username, password)
	if err != nil {
		return err
	}
	data := MongoBindingCredentials{
		Host:             e.Host,
		Port:             e.Port,
		Username:         username,
		Password:         password,
		ConnectionString: generateConnectionString(e.Host, e.Port, request.InstanceID, username, password),
	}

	*response = data
	return nil
}

func (e *MongoDriver) CredentialsExist(request model.CredentialsRequest, response *bool) error {
	username := request.InstanceID + "-" + request.CredentialsID

	created, err := e.db.IsUserCreated(request.InstanceID, username)
	if err != nil {
		return err
	}
	*response = created
	return nil
}

func (e *MongoDriver) RevokeCredentials(request model.CredentialsRequest, response *interface{}) error {
	username := request.InstanceID + "-" + request.CredentialsID
	err := e.db.DeleteUser(request.InstanceID, username)
	if err != nil {
		return err
	}
	*response = fmt.Sprintf("Credentials for %s revoked", username)
	return nil
}

func (e *MongoDriver) DeprovisionInstance(instanceID string, response *interface{}) error {
	err := e.db.DeleteDatabase(instanceID)
	if err != nil {
		return err
	}
	*response = fmt.Sprintf("Deprovisioned %s", instanceID)

	return nil
}
