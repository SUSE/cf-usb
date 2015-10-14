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

func Test_Bind(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("CreateUser", "testId", "testId-user", mock.Anything).Return(nil)
	driver.db = mockProv

	var req model.CredentialsRequest

	req.InstanceID = "testId"
	req.CredentialsID = "user"

	var response interface{}
	err := driver.GenerateCredentials(req, &response)

	assert.NoError(err)
}

func Test_Unbind(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)

	mockProv.On("DeleteUser", "testId-user").Return(nil)
	driver.db = mockProv

	var req model.CredentialsRequest

	req.InstanceID = "testId"
	req.CredentialsID = "user"
	var response interface{}
	err := driver.RevokeCredentials(req, &response)

	assert.NoError(err)
}
