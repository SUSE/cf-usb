package driver

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type MongoDriver struct {
	User   string `json:"user id"`
	Pass   string `json:"password"`
	Host   string `json:"server"`
	Port   string `json:"port"`
	db     mongoprovisioner.MongoProvisionerInterface
	logger lager.Logger
	driver.Driver
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

func (e *MongoDriver) Init(configuration config.DriverProperties, response *string) error {
	err := json.Unmarshal(*configuration.DriverConfiguration, &e)
	e.logger.Info("Mysql Driver initializing")
	e.db, err = mongoprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port, e.logger)
	return err
}

func (e *MongoDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	created, err := e.db.IsDatabaseCreated(request.InstanceID)

	if err != nil {
		return err
	}

	if created {
		return brokerapi.ErrInstanceAlreadyExists
	} else {
		err := e.db.CreateDatabase(request.InstanceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *MongoDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	created, err := e.db.IsDatabaseCreated(request.InstanceID)

	if err != nil {
		return err
	}
	if created {
		err := e.db.DeleteDatabase(request.InstanceID)
		if err != nil {
			return err
		}
		*response = fmt.Sprintf("Deprovisioned %s", request.InstanceID)
	} else {
		return brokerapi.ErrInstanceDoesNotExist
	}
	return nil
}

func (e *MongoDriver) Bind(request model.DriverBindRequest, response *json.RawMessage) error {

	username := request.InstanceID + "-" + request.BindingID
	password := e.secureRandomString(32)

	created, err := e.db.IsUserCreated(request.InstanceID, username)
	if err != nil {
		return err
	}

	if created {
		return brokerapi.ErrBindingAlreadyExists
	} else {

		err := e.db.CreateUser(request.InstanceID, username, password)
		if err != nil {
			return err
		}
		//TO DO: Add response object
	}
	return nil
}

func (e *MongoDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	username := request.InstanceID + "-" + request.BindingID
	created, err := e.db.IsUserCreated(request.InstanceID, username)
	if err != nil {
		return err
	}
	if created {
		err := e.db.DeleteUser(request.InstanceID, username)
		if err != nil {
			return err
		}
	} else {
		return brokerapi.ErrBindingDoesNotExist
	}
	*response = fmt.Sprintf("Unbound %s", username)
	return nil
}
