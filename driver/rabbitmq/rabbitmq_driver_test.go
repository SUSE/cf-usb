package rabbitmq

import (
	"encoding/json"
	"errors"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/config"
	"github.com/hpcloud/cf-usb/driver/rabbitmq/rabbitmqprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("rabbitmq-driver-test")

func getEmptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	emptyConfig := rawMessage
	return &emptyConfig
}

func getMockProvisioner() (*mocks.RabbitmqProvisionerInterface, usbDriver.Driver) {
	mockProv := new(mocks.RabbitmqProvisionerInterface)
	rabbitDriver := NewRabbitmqDriver(logger, mockProv)

	mockProv.On("Connect", config.RabbitmqDriverConfig{}).Return(nil)

	return mockProv, rabbitDriver
}

func Test_Ping(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()

	mockProv.On("PingServer").Return(nil)

	var response bool
	err := rabbitDriver.Ping(getEmptyConfig(), &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_PingFail(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("PingServer").Return(errors.New("ping fail"))

	var response bool
	err := rabbitDriver.Ping(getEmptyConfig(), &response)

	assert.Error(err)
	assert.False(response)
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("CreateContainer", "testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	req.Dials = getEmptyConfig()

	var response usbDriver.Instance
	err := rabbitDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("ContainerExists", "testId").Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Instance

	err := rabbitDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	_, rabbitDriver := getMockProvisioner()

	var response string
	err := rabbitDriver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	_, rabbitDriver := getMockProvisioner()

	var response string
	err := rabbitDriver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("UserExists", "testId", "user").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := rabbitDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetCredentialsNotExist(t *testing.T) {
	assert := assert.New(t)
	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("UserExists", "testId", "fakeUser").Return(false, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "fakeUser"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := rabbitDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.DoesNotExist)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("CreateUser", mock.Anything, mock.Anything).Return(map[string]string{
		"host":     "127.0.0.1",
		"user":     "user",
		"password": "password",
		"port":     "1234",
	}, nil)

	var req usbDriver.GenerateCredentialsRequest
	req.InstanceID = "testContainer"
	req.Config = getEmptyConfig()
	var response interface{}

	err := rabbitDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("DeleteUser", "testId", "testIduser").Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "testIduser"
	req.InstanceID = "testId"
	var response usbDriver.Credentials

	err := rabbitDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, rabbitDriver := getMockProvisioner()
	mockProv.On("DeleteContainer", "testContainer").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEmptyConfig()
	req.InstanceID = "testContainer"

	var response usbDriver.Instance
	err := rabbitDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
