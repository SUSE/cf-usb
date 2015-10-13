package postgres

import (
	"encoding/json"
	"testing"

	"github.com/hpcloud/cf-usb/driver/postgres/postgresprovisioner"
	"github.com/hpcloud/cf-usb/lib/model"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner")

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("CreateDatabase", "testId").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(false, nil)
	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverProvisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Provision(req, &response)

	assert.NoError(err)
}

func Test_ProvisionExists(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("CreateDatabase", "testId").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverProvisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Provision(req, &response)

	assert.Error(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverDeprovisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Deprovision(req, &response)

	assert.NoError(err)
}

func Test_DeprovisionNotExists(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(false, nil)
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverDeprovisionRequest

	req.InstanceID = "testId"

	var response string
	err := driver.Deprovision(req, &response)

	assert.Error(err)
}

func Test_Bind(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	mockProv.On("UserExists", "testIduser").Return(false, nil)
	mockProv.On("CreateUser", "testId", "testIduser", mock.Anything).Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverBindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"

	var response json.RawMessage
	err := driver.Bind(req, &response)

	assert.NoError(err)
}

func Test_BindNoDb(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(false, nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverBindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"

	var response json.RawMessage
	err := driver.Bind(req, &response)

	assert.Error(err)
}

func Test_BindExists(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	mockProv.On("UserExists", "testIduser").Return(true, nil)
	mockProv.On("CreateUser", "testId", "testIduser", mock.Anything).Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverBindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"

	var response json.RawMessage
	err := driver.Bind(req, &response)

	assert.Error(err)
}

func Test_Unbind(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	mockProv.On("UserExists", "testIduser").Return(true, nil)
	mockProv.On("DeleteUser", "testId", "testIduser").Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverUnbindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"
	var response string
	err := driver.Unbind(req, &response)

	assert.NoError(err)
}

func Test_UnbindNoDb(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(false, nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverUnbindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"
	var response string
	err := driver.Unbind(req, &response)

	assert.Error(err)
}

func Test_UnbindNotExists(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DatabaseExists", "testId").Return(true, nil)
	mockProv.On("UserExists", "testIduser").Return(false, nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.DriverUnbindRequest

	req.InstanceID = "testId"
	req.BindingID = "user"
	var response string
	err := driver.Unbind(req, &response)

	assert.Error(err)
}
