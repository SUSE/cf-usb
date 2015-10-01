package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hpcloud/cf-usb/driver/postgresprovisioner"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/hpcloud/gocfbroker"
)

type PostgresServiceProperties struct {
	DefaultPostgresConnection map[string]string `json:"defaultPostgresConnection"`
}

type postgresDriver struct {
	driverProperties config.DriverProperties
	driverConnParams PostgresServiceProperties
	Driver
}

func NewPostgresDriver() Driver {
	return &postgresDriver{}
}

func (driver *postgresDriver) Init(driverProperties config.DriverProperties, response *string) error {
	driver.driverProperties = driverProperties

	log.Println("Driver initialized")
	for _, service := range driverProperties.Services {
		log.Println("Using serviceID:", service.ID)
		log.Println("Service Description", service.Description)
		for _, plan := range service.Plans {
			log.Println("PlanID:", plan.ID)
			log.Println("PlanName:", plan.Name)
		}
	}

	conf := (*json.RawMessage)(driverProperties.DriverConfiguration)
	log.Println(string(*conf))
	dsp := PostgresServiceProperties{}
	err := json.Unmarshal(*conf, &dsp)
	if err != nil {
		return err
	}
	driver.driverConnParams = dsp

	*response = "Sucessfully initialized postgres driver"
	return nil
}

func (driver *postgresDriver) Provision(request model.DriverProvisionRequest, response *string) error {
	postgresprovisioner := postgresprovisioner.NewPostgresProvisioner(driver.driverConnParams.DefaultPostgresConnection)

	err := postgresprovisioner.CreateDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	*response = fmt.Sprintf("http://example-dashboard.com/9189kdfsk0vfnku")
	return nil
}

func (driver *postgresDriver) Deprovision(request model.DriverDeprovisionRequest, response *string) error {
	postgresprovisioner := postgresprovisioner.NewPostgresProvisioner(driver.driverConnParams.DefaultPostgresConnection)

	err := postgresprovisioner.DeleteDatabase(request.InstanceID)
	if err != nil {
		return err
	}

	*response = "Successfully deprovisoned"
	return nil
}

func (driver *postgresDriver) Update(request model.DriverUpdateRequest, response *string) error {

	return nil
}

func (driver *postgresDriver) Bind(request model.DriverBindRequest, response *gocfbroker.BindingResponse) error {
	username := request.InstanceID + request.BindingID
	password := secureRandomString(32)

	postgresprovisioner := postgresprovisioner.NewPostgresProvisioner(driver.driverConnParams.DefaultPostgresConnection)

	err := postgresprovisioner.CreateUser(request.InstanceID, username, password)
	if err != nil {
		return err
	}

	data := []byte(fmt.Sprintf(`{"username": "%v", "password": "%v"}`, username, password))
	response.Credentials = (*json.RawMessage)(&data)

	return nil
}

func (driver *postgresDriver) Unbind(request model.DriverUnbindRequest, response *string) error {
	username := request.InstanceID + request.BindingID

	postgresprovisioner := postgresprovisioner.NewPostgresProvisioner(driver.driverConnParams.DefaultPostgresConnection)

	err := postgresprovisioner.DeleteUser(request.InstanceID, username)
	if err != nil {
		return err
	}

	return nil
}

func secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		log.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}
