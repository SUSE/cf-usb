package postgres

import (
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
	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.ProvisionInstanceRequest

	req.InstanceID = "testId"

	var response bool
	err := driver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DeleteDatabase", "testId").Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	req := "testId"

	var response interface{}
	err := driver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}

func Test_Bind(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("CreateUser", "testId", "testIduser", mock.Anything).Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.CredentialsRequest

	req.InstanceID = "testId"
	req.CredentialsID = "user"

	var response interface{}
	err := driver.GenerateCredentials(req, &response)

	assert.NoError(err)
}

func Test_Unbind(t *testing.T) {
	assert := assert.New(t)

	driver := postgresDriver{}
	driver.logger = logger
	mockProv := new(postgresprovisioner.PostgresProvisionerMock)
	mockProv.On("Init").Return(nil)
	mockProv.On("DeleteUser", "testId", "testIduser").Return(nil)

	driver.postgresProvisioner = mockProv

	driver.postgresProvisioner.Init()

	var req model.CredentialsRequest

	req.InstanceID = "testId"
	req.CredentialsID = "user"
	var response interface{}
	err := driver.RevokeCredentials(req, &response)

	assert.NoError(err)
}
