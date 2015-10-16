package redis

import (
	"errors"
	"testing"

	"github.com/hpcloud/cf-usb/driver/redis/redisprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-provisioner")

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("CreateContainer", "testContainer").Return(nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var req model.ProvisionInstanceRequest
	req.InstanceID = "testContainer"

	var response bool
	err := driver.ProvisionInstance(req, &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DeleteContainer", "testContainer").Return(nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	req := "testContainer"

	var response interface{}
	err := driver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}

func Test_InstanceExists(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("ContainerExists", "testContainer").Return(true, nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	req := "testContainer"

	var response bool
	err := driver.InstanceExists(req, &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_InstanceNotExists(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("ContainerExists", "testContainer").Return(false, nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	req := "testContainer"

	var response bool
	err := driver.InstanceExists(req, &response)

	assert.NoError(err)
	assert.False(response)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	var response string
	err := driver.GetDailsSchema("", &response)

	assert.NoError(err)
	assert.NotNil(response)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	var response string
	err := driver.GetConfigSchema("", &response)

	assert.NoError(err)
	assert.NotNil(response)
}

func Test_CredentialsExist(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("ContainerExists", "testContainer").Return(true, nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var req model.CredentialsRequest
	req.InstanceID = "testContainer"

	var response bool
	err := driver.CredentialsExist(req, &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_CredentialNotExist(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("ContainerExists", "testContainer").Return(false, nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var req model.CredentialsRequest
	req.InstanceID = "testContainer"

	var response bool
	err := driver.CredentialsExist(req, &response)

	assert.NoError(err)
	assert.False(response)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("GetCredentials", "testContainer").Return(map[string]string{
		"password": "password",
		"port":     "1234",
	}, nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var req model.CredentialsRequest
	req.InstanceID = "testContainer"

	var response interface{}
	err := driver.GenerateCredentials(req, &response)

	assert.NoError(err)
	assert.NotNil(response)

	credentials := response.(RedisBindingCredentials)
	assert.Equal(credentials.Password, "password")
	assert.Equal(credentials.Port, "1234")
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	var req model.CredentialsRequest
	req.InstanceID = "testContainer"

	var response interface{}
	err := driver.RevokeCredentials(req, &response)

	assert.NoError(err)
}

func Test_Ping(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("PingServer").Return(nil)
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var response bool
	err := driver.Ping("", &response)

	assert.NoError(err)
	assert.True(response)
}

func Test_PingFail(t *testing.T) {
	assert := assert.New(t)

	driver := redisDriver{}
	driver.logger = logger

	mockProv := new(redisprovisioner.RedisProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("PingServer").Return(errors.New("ping fail"))
	driver.redisProvisioner = mockProv

	driver.redisProvisioner.Init()

	var response bool
	err := driver.Ping("", &response)

	assert.Error(err)
	assert.False(response)
}
