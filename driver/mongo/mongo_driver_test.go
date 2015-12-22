package mongo

import (
	"encoding/json"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mongo/config"
	"github.com/hpcloud/cf-usb/driver/mongo/mongoprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mongo-driver-test")

func getEmptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	emptyConfig := rawMessage
	return &emptyConfig
}

func getMockProvisioner() (*mocks.MongoProvisionerInterface, usbDriver.Driver) {
	mockProv := new(mocks.MongoProvisionerInterface)
	mongoDriver := NewMongoDriver(logger, mockProv)

	mockProv.On("Connect", config.MongoDriverConfig{}).Return(nil)

	return mockProv, mongoDriver
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("CreateDatabase", "testId").Return(nil)
	var req usbDriver.ProvisionInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	req.Dials = getEmptyConfig()
	var response usbDriver.Instance

	err := mongoDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("IsDatabaseCreated", "testId").Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Instance

	err := mongoDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	_, mongoDriver := getMockProvisioner()

	var response string
	err := mongoDriver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	_, mongoDriver := getMockProvisioner()

	var response string
	err := mongoDriver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("IsUserCreated", "testId", "testId-user").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := mongoDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("CreateUser", "testId", "testId-user", mock.Anything).Return(nil)

	var req usbDriver.GenerateCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response interface{}

	err := mongoDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("DeleteUser", "testId", "testId-user").Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Credentials

	err := mongoDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mongoDriver := getMockProvisioner()
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEmptyConfig()
	req.InstanceID = "testId"

	var response usbDriver.Instance
	err := mongoDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
