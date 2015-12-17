package redis

import (
	"encoding/json"
	"errors"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/redis/config"
	"github.com/hpcloud/cf-usb/driver/redis/redisprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-provisioner")

func getEmptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	emptyConfig := rawMessage
	return &emptyConfig
}

func getMockProvisioner() (*mocks.RedisProvisionerInterface, usbDriver.Driver) {
	mockProv := new(mocks.RedisProvisionerInterface)
	redisDriver := NewRedisDriver(logger, mockProv)

	mockProv.On("Connect", config.RedisDriverConfig{}).Return(nil)

	return mockProv, redisDriver
}

func Test_Ping(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()

	mockProv.On("PingServer").Return(nil)

	var response bool
	err := redisDriver.Ping(getEmptyConfig(), &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_PingFail(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("PingServer").Return(errors.New("ping fail"))

	var response bool
	err := redisDriver.Ping(getEmptyConfig(), &response)

	assert.Error(err)
	assert.False(response)
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("CreateContainer", "testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	req.Dials = getEmptyConfig()

	var response usbDriver.Instance
	err := redisDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("ContainerExists", mock.Anything).Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Instance

	err := redisDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	_, redisDriver := getMockProvisioner()

	var response string
	err := redisDriver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	_, redisDriver := getMockProvisioner()

	var response string
	err := redisDriver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("UserExists", "testIduser").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := redisDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.DoesNotExist)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("GetCredentials", "testContainer").Return(map[string]string{
		"password": "password",
		"port":     "1234",
	}, nil)

	var req usbDriver.GenerateCredentialsRequest
	req.InstanceID = "testContainer"
	req.Config = getEmptyConfig()
	var response interface{}

	err := redisDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("DeleteUser", "testId", "testIduser").Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.Config = getEmptyConfig()
	var response usbDriver.Credentials

	err := redisDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, redisDriver := getMockProvisioner()
	mockProv.On("DeleteContainer", "testContainer").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEmptyConfig()
	req.InstanceID = "testContainer"

	var response usbDriver.Instance
	err := redisDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
