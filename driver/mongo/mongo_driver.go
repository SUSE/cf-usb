package mongo

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mongo/config"
	"github.com/hpcloud/cf-usb/driver/mongo/driverdata"
	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type MongoDriver struct {
	logger lager.Logger
	conf   config.MongoDriverConfig
	db     mongoprovisioner.MongoProvisionerInterface
}

func (e *MongoDriver) secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		e.logger.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func NewMongoDriver(logger lager.Logger, db mongoprovisioner.MongoProvisionerInterface) driver.Driver {
	return &MongoDriver{logger: logger, db: db}
}

func (e *MongoDriver) init(configuration *json.RawMessage) error {
	mongoConfig := config.MongoDriverConfig{}
	err := json.Unmarshal(*configuration, &mongoConfig)
	if err != nil {
		return err
	}
	e.logger.Info("Mongo Driver initializing")
	err = e.db.Connect(mongoConfig)
	e.conf = mongoConfig
	return err
}

func (e *MongoDriver) Ping(request *json.RawMessage, result *bool) error {
	err := e.init(request)
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

func (e *MongoDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	err = e.db.CreateDatabase(request.InstanceID)
	if err != nil {
		return err
	}
	response.Status = status.Created
	return nil
}

func (e *MongoDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	created, err := e.db.IsDatabaseCreated(request.InstanceID)
	if err != nil {
		return err
	}
	response.Status = status.DoesNotExist
	if created {
		response.Status = status.Exists
	}
	return nil
}

func (e *MongoDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	username := request.InstanceID + "-" + request.CredentialsID
	password := e.secureRandomString(32)

	err = e.db.CreateUser(request.InstanceID, username, password)
	if err != nil {
		return err
	}
	data := MongoBindingCredentials{
		Host:             e.conf.Host,
		Port:             e.conf.Port,
		Username:         username,
		Password:         password,
		ConnectionString: generateConnectionString(e.conf.Host, e.conf.Port, request.InstanceID, username, password),
	}

	*response = data
	return nil
}

func (e *MongoDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	username := request.InstanceID + "-" + request.CredentialsID

	created, err := e.db.IsUserCreated(request.InstanceID, username)
	if err != nil {
		return err
	}
	response.Status = status.DoesNotExist
	if created {
		response.Status = status.Exists
	}
	return nil
}

func (e *MongoDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	username := request.InstanceID + "-" + request.CredentialsID

	err = e.db.DeleteUser(request.InstanceID, username)
	if err != nil {
		return err
	}
	response.Status = status.Deleted
	return nil
}

func (e *MongoDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	err = e.db.DeleteDatabase(request.InstanceID)
	if err != nil {
		return err
	}
	response.Status = status.Deleted

	return nil
}
