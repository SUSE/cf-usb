package driver

import (
	"encoding/json"
	"os"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mysql/config"
	"github.com/hpcloud/cf-usb/driver/mysql/mysqlprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getEnptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	emptyConfig := rawMessage
	return &emptyConfig
}

func getMockProvisioner() (*mocks.MysqlProvisionerInterface, usbDriver.Driver) {
	var logger = lager.NewLogger("mongodb-driver-test")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	mockProv := new(mocks.MysqlProvisionerInterface)
	mysqlDriver := NewMysqlDriver(logger, mockProv)

	mockProv.On("Connect", config.MysqlDriverConfig{}).Return(nil)

	return mockProv, mysqlDriver
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("CreateDatabase", "testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getEnptyConfig()

	var response usbDriver.Instance
	err := mysqlDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("IsDatabaseCreated", mock.Anything).Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response usbDriver.Instance

	err := mysqlDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	driver := MysqlDriver{}

	var response string
	err := driver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	driver := MysqlDriver{}

	var response string
	err := driver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("IsUserCreated", "ee11cbb19052e40b").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEnptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := mysqlDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("CreateUser", "testId", mock.Anything, mock.Anything).Return(nil)

	var req usbDriver.GenerateCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response interface{}

	err := mysqlDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("DeleteUser", mock.Anything).Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response usbDriver.Credentials

	err := mysqlDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mysqlDriver := getMockProvisioner()
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEnptyConfig()
	req.InstanceID = "testId"

	var response usbDriver.Instance
	err := mysqlDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
