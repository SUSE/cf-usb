package driver

import (
	"encoding/json"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/mssql/config"
	"github.com/hpcloud/cf-usb/driver/mssql/mssqlprovisioner/mocks"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mssql-driver-test")

func getEmptyConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte("{}"))
	return &rawMessage
}

func getConfigMessage(conf config.MssqlDriverConfig) *json.RawMessage {
	raw, _ := json.Marshal(conf)
	ret := json.RawMessage(raw)
	return &ret
}

func getMockProvisioner() (*mocks.MssqlProvisionerInterface, usbDriver.Driver) {
	mockProv := new(mocks.MssqlProvisionerInterface)
	mssqlDriver := NewMssqlDriver(logger, mockProv)

	mockProv.On("Connect", "mssql", mock.Anything).Return(nil)

	return mockProv, mssqlDriver
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("CreateDatabase", "testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	req.Dials = getEmptyConfig()

	var response usbDriver.Instance
	err := mssqlDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_ProvisionWithPrefix(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("CreateDatabase", "cf-testId").Return(nil)

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getConfigMessage(config.MssqlDriverConfig{
		DbIdentifierPrefix: "cf-",
	})
	req.Dials = getEmptyConfig()

	var response usbDriver.Instance
	err := mssqlDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("IsDatabaseCreated", mock.Anything).Return(true, nil)
	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Instance

	err := mssqlDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	_, mssqlDriver := getMockProvisioner()

	var response string
	err := mssqlDriver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	_, mssqlDriver := getMockProvisioner()

	var response string
	err := mssqlDriver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("IsUserCreated", "testId", "user").Return(true, nil)

	var req usbDriver.GetCredentialsRequest
	req.Config = getEmptyConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := mssqlDriver.GetCredentials(req, &response)
	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("CreateUser", "testId", mock.Anything, mock.Anything).Return(nil)

	var req usbDriver.GenerateCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response interface{}

	err := mssqlDriver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("DeleteUser", "testId", "user").Return(nil)

	var req usbDriver.RevokeCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getEmptyConfig()
	var response usbDriver.Credentials

	err := mssqlDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	mockProv, mssqlDriver := getMockProvisioner()
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getEmptyConfig()
	req.InstanceID = "testId"

	var response usbDriver.Instance
	err := mssqlDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
