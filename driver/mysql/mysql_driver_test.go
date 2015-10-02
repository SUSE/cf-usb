package driver

import (
	"testing"

	"encoding/json"
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

	var req model.DriverProvisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Provision(req, &response)

	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("DeleteDatabase", "testId").Return(nil)
	driver.db = mockProv

	var req model.DriverDeprovisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Deprovision(req, &response)

	assert.NoError(err)
}

func Test_Bind(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("CreateUser", "testId", "testId-user", mock.Anything).Return(nil)
	driver.db = mockProv

	var req model.DriverBindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"

	var response json.RawMessage
	err := driver.Bind(req, &response)

	assert.NoError(err)
}

func Test_Unbind(t *testing.T) {
	assert := assert.New(t)

	driver := MysqlDriver{}
	mockProv := new(mysqlprovisioner.MysqlProvisionerMock)
	mockProv.On("DeleteUser", "testId-user").Return(nil)
	driver.db = mockProv

	var req model.DriverUnbindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"
	var response string
	err := driver.Unbind(req, &response)

	assert.NoError(err)
}
