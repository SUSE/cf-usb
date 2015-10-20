package driver

import (
	"testing"

	"github.com/hpcloud/cf-usb/driver/mysql/mysqlprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("CreateDatabase", "testId").Return(nil)
	driver.db = mockProv

	var req model.ProvisionInstanceRequest

	req.InstanceID = "testId"

	var response bool
	err := driver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_InstanceExists(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("IsDatabaseCreated", "testId").Return(true, nil)
	driver.db = mockProv
	req := "testId"

	var response bool
	err := driver.InstanceExists(req, &response)

	assert.True(response)
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

func Test_CredentialsExist(t *testing.T) {
	assert := assert.New(t)
	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("IsUserCreated", mock.Anything).Return(true, nil)

	driver.db = mockProv

	var req model.CredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response bool
	err := driver.CredentialsExist(req, &response)
	assert.True(response)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)
	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("CreateUser", "testId", mock.Anything, mock.Anything).Return(nil)

	driver.db = mockProv

	var req model.CredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	var response interface{}

	err := driver.GenerateCredentials(req, &response)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("DeleteUser", mock.Anything).Return(nil)

	driver.db = mockProv

	var req model.CredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	var response interface{}

	err := driver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	driver.db = mockProv

	req := "testId"

	var response interface{}
	err := driver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
