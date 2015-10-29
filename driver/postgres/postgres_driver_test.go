package postgres

import (
	"encoding/json"
	"os"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/postgres/config"
	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner")

func getEnptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	emptyConfig := rawMessage
	return &emptyConfig
}

func getMockProvisioner() (*mocks.PostgresProvisionerInterface, usbDriver.Driver) {
	var logger = lager.NewLogger("mongodb-driver-test")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	mockProv := new(mocks.PostgresProvisionerInterface)
	postgresDriver := NewPostgresDriver(logger, mockProv)

	mockProv.On("Connect", config.PostgresDriverConfig{}).Return(nil)

	return mockProv, postgresDriver
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("CreateDatabase", "testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getEnptyConfig()

	var response usbDriver.Instance
	err := postgresDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("DatabaseExists", mock.Anything).Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response usbDriver.Instance

	err := postgresDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	driver := PostgresDriver{}

	var response string
	err := driver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	driver := PostgresDriver{}

	var response string
	err := driver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("UserExists", "testIduser").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEnptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := postgresDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("CreateUser", "testId", mock.Anything, mock.Anything).Return(nil)

	var req usbDriver.GenerateCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response interface{}

	err := postgresDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("DeleteUser", "testId", "testIduser").Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEnptyConfig()
	var response usbDriver.Credentials

	err := postgresDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, postgresDriver := getMockProvisioner()
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEnptyConfig()
	req.InstanceID = "testId"

	var response usbDriver.Instance
	err := postgresDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
