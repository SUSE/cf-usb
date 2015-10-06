package postgresdriver

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	//"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

var postgresProvisioner *postgresprovisioner.PostgresProvisioner

type postgresDriver struct {
	driverProperties config.DriverProperties
	logger           lager.Logger
	driver.Driver
}

func NewPostgresDriver(logger lager.Logger) driver.Driver {
	return &postgresDriver{logger: logger}
}

func (driver *postgresDriver) Init(driverProperties config.DriverProperties, response *string) error {
	driver.driverProperties = driverProperties

	log.Println("Driver initialized")

	conf := (*json.RawMessage)(driverProperties.DriverConfiguration)
	log.Println(string(*conf))
	dsp := postgresprovisioner.PostgresServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	driver.logger.Info("init-driver", lager.Data{"user": dsp.User, "password": dsp.Password})

	postgresProvisioner := postgresprovisioner.NewPostgresProvisioner(dsp, driver.logger)
	postgresProvisioner.Init()
	if err != nil {
		driver.logger.Fatal("error-initializing-provisioner", err)
		return err
	}

	*response = "Sucessfully initialized postgres driver"
	return nil
}

func (driver *postgresDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	driver.logger.Info("Provisioning", lager.Data{"instance-id": request.InstanceID, "plan-id": request.ServiceDetails.PlanID})

	err := postgresProvisioner.CreateDatabase(request.InstanceID)
	if err != nil {
		return err
	}
	*response = ""
	return nil
}

func (driver *postgresDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	driver.logger.Info("deprovision-request", lager.Data{"instance-id": request.InstanceID})

	err := postgresProvisioner.DeleteDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	*response = "Successfully deprovisoned"
	return nil
}

func (driver *postgresDriver) Bind(request model.DriverBindRequest, response *json.RawMessage) error {
	driver.logger.Info("bind-request", lager.Data{"instanceID": request.InstanceID,
		"planID": request.BindDetails.PlanID, "appID": request.BindDetails.AppGUID})

	username := request.InstanceID + request.BindingID
	password, err := secureRandomString(32)
	if err != nil {
		return err
	}

	err = postgresProvisioner.CreateUser(request.InstanceID, username, password)
	if err != nil {
		return err
	}

	data := []byte(fmt.Sprintf(`{"username": "%v", "password": "%v"}`, username, password))
	response = (*json.RawMessage)(&data)

	return nil
}

func (driver *postgresDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	driver.logger.Info("unbind-request", lager.Data{"bindingID": request.BindingID, "InstanceID": request.InstanceID})
	username := request.InstanceID + request.BindingID

	err := postgresProvisioner.DeleteUser(request.InstanceID, username)
	if err != nil {
		return err
	}

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
