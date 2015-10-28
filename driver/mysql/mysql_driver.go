package driver

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mysql/config"
	"github.com/hpcloud/cf-usb/driver/mysql/driverdata"
	"github.com/hpcloud/cf-usb/driver/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
)

type MysqlDriver struct {
	logger lager.Logger
	conf   config.MysqlDriverConfig
	db     mysqlprovisioner.MysqlProvisionerInterface
}

func (e *MysqlDriver) secureRandomString(bytesOfEntpry int) string {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		e.logger.Fatal("rng-failure", err)
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func (e *MysqlDriver) getMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

func NewMysqlDriver(logger lager.Logger, db mysqlprovisioner.MysqlProvisionerInterface) driver.Driver {
	return &MysqlDriver{logger: logger, db: db}
}

func (e *MysqlDriver) init(conf *json.RawMessage) error {
	mysqlConfig := config.MysqlDriverConfig{}
	err := json.Unmarshal(*conf, &mysqlConfig)
	e.logger.Info("Mysql Driver initializing")
	err = e.db.Connect(mysqlConfig)
	if err != nil {
		return err
	}
	e.conf = mysqlConfig
	return nil
}

func (e *MysqlDriver) Ping(request *json.RawMessage, response *bool) error {
	*response = false

	err := e.init(request)
	if err != nil {
		return err
	}
	*response = true
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
func (e *MysqlDriver) ProvisionInstance(request driver.ProvisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	err = e.db.CreateDatabase(strings.Replace(request.InstanceID, "-", "", -1))
	if err != nil {
		return err
	}
	response.Status = status.Created
	return nil
}

func (e *MysqlDriver) GetInstance(request driver.GetInstanceRequest, response *driver.Instance) error {
	response.Status = status.DoesNotExist
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	created, err := e.db.IsDatabaseCreated(strings.Replace(request.InstanceID, "-", "", -1))
	if err != nil {
		return err
	}

	if created {
		response.Status = status.Exists
	}

	return nil
}

func (e *MysqlDriver) GenerateCredentials(request driver.GenerateCredentialsRequest, response *interface{}) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}
	username, err := e.getMD5Hash(request.CredentialsID)
	if err != nil {
		return err
	}
	if len(username) > 16 {
		username = username[:16]
	}
	password := e.secureRandomString(32)

	err = e.db.CreateUser(strings.Replace(request.InstanceID, "-", "", -1), username, password)
	if err != nil {
		return err
	}
	data := MysqlBindingCredentials{
		Host:     e.conf.Host,
		Port:     e.conf.Port,
		Username: username,
		Password: password,
		ConnectionString: generateConnectionString(e.conf.Host, e.conf.Port,
			strings.Replace(request.InstanceID, "-", "", -1), username, password),
	}

	*response = data
	return nil
}

func (e *MysqlDriver) GetCredentials(request driver.GetCredentialsRequest, response *driver.Credentials) error {
	response.Status = status.DoesNotExist
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	username, err := e.getMD5Hash(request.CredentialsID)
	if err != nil {
		return err
	}
	if len(username) > 16 {
		username = username[:16]
	}
	created, err := e.db.IsUserCreated(username)
	if err != nil {
		return err
	}
	if created {
		response.Status = status.Exists
	}
	return nil
}

func (e *MysqlDriver) RevokeCredentials(request driver.RevokeCredentialsRequest, response *driver.Credentials) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	username, err := e.getMD5Hash(request.CredentialsID)
	if err != nil {
		return err
	}
	if len(username) > 16 {
		username = username[:16]
	}
	err = e.db.DeleteUser(username)
	if err != nil {
		return err
	}

	response.Status = status.Deleted
	return nil
}

func (e *MysqlDriver) DeprovisionInstance(request driver.DeprovisionInstanceRequest, response *driver.Instance) error {
	err := e.init(request.Config)
	if err != nil {
		return err
	}

	err = e.db.DeleteDatabase(strings.Replace(request.InstanceID, "-", "", -1))
	if err != nil {
		return err
	}
	response.Status = status.Deleted

	return nil
}
