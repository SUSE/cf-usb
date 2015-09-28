package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hpcloud/cf-usb/driver/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/hpcloud/gocfbroker"
)

type MysqlDriver struct {
	User string `json:"user id"`
	Pass string `json:"password"`
	Host string `json:"server"`
	Port string `json:"port"`
	Driver
}

func secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		log.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func NewMysqlDriver() Driver {
	return &MysqlDriver{}
}

func (e *MysqlDriver) Init(configuration config.DriverProperties, response *string) error {
	err := json.Unmarshal(*configuration.DriverConfiguration, &e)
	log.Println("Mysql Driver initializing")
	return err
}

func (e *MysqlDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	db := mysqlprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port)

	err := db.CreateDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	*response = fmt.Sprintf("http://localhost/instance/%s/service/%s/dashboard", request.InstanceID, request.BrokerProvisionRequest.ServiceID)

	return nil
}

func (e *MysqlDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	db := mysqlprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port)

	err := db.DeleteDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	*response = fmt.Sprintf("Deprovisioned %s", request.InstanceID)

	return nil
}

func (e *MysqlDriver) Update(request model.DriverUpdateRequest, response *string) error {
	log.Printf("\n\nUpdate called with:\ninstanceID: %s", request.InstanceID)
	return nil
}

func (e *MysqlDriver) Bind(request model.DriverBindRequest, response *gocfbroker.BindingResponse) error {
	username := request.InstanceID + "-user"
	password := secureRandomString(32)

	log.Println("BINDING ...")
	log.Println(username)
	log.Println(password)

	db := mysqlprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port)

	err := db.CreateUser(request.InstanceID, username, password)
	if err != nil {
		return err
	}
	data := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	response.Credentials = (*json.RawMessage)(&data)

	return nil
}

func (e *MysqlDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	username := request.InstanceID + "-user"

	db := mysqlprovisioner.New(e.User, e.Pass, e.Host+":"+e.Port)

	err := db.DeleteUser(username)
	if err != nil {
		return err
	}

	return nil
}

func (e *MysqlDriver) GetCatalog(request string, response *string) error {
	return nil
}

func (e *MysqlDriver) GetInstances(request string, response *string) error {
	return nil
}
