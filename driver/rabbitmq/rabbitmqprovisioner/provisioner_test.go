package rabbitmqprovisioner

import (
	"github.com/hpcloud/cf-usb/driver/rabbitmq/config"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("rabbitmq-provisioner")

var testRabbitmqProv = struct {
	rabbitmqProvisioner RabbitmqProvisionerInterface
	driverConfig        config.RabbitmqDriverConfig
}{}

func init() {
	testRabbitmqProv.driverConfig = config.RabbitmqDriverConfig{
		DockerEndpoint: os.Getenv("DOCKER_ENDPOINT"),
		DockerImage:    os.Getenv("RABBIT_DOCKER_IMAGE"),
		ImageVersion:   os.Getenv("RABBIT_DOCKER_IMAGE_VERSION"),
	}

	testRabbitmqProv.rabbitmqProvisioner = NewRabbitmqProvisioner(logger)
	testRabbitmqProv.rabbitmqProvisioner.Connect(testRabbitmqProv.driverConfig)
}

func TestCreateContainer(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"

	err := testRabbitmqProv.rabbitmqProvisioner.CreateContainer(name)
	assert.NoError(err)

	exists, err := testRabbitmqProv.rabbitmqProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.True(exists)
}

func TestCreateUser(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	time.Sleep(9 * time.Second)
	assert := assert.New(t)

	name := "rabbitContainer"
	credentialId := "someCred"

	credentials, err := testRabbitmqProv.rabbitmqProvisioner.CreateUser(name, credentialId)
	assert.NoError(err)
	assert.NotNil(credentials["password"])
	assert.NotNil(credentials["port"])
	assert.NotNil(credentials["mgmt_port"])
	assert.NotNil(credentials["host"])
	assert.NotNil(credentials["user"])
	assert.NotNil(credentials["vhost"])
}

func TestUserExists(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"
	credentialId := "someCred"

	exists, err := testRabbitmqProv.rabbitmqProvisioner.UserExists(name, credentialId)
	assert.NoError(err)
	assert.True(exists)
}

func TestDeleteUser(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"
	credentialId := "someCred"

	err := testRabbitmqProv.rabbitmqProvisioner.DeleteUser(name, credentialId)
	assert.NoError(err)
}

func TestUserNotExists(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"
	credentialId := "someCred"

	exists, err := testRabbitmqProv.rabbitmqProvisioner.UserExists(name, credentialId)
	assert.NoError(err)
	assert.False(exists)
}

func TestDeleteContainer(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"

	err := testRabbitmqProv.rabbitmqProvisioner.DeleteContainer(name)
	assert.NoError(err)

	exists, err := testRabbitmqProv.rabbitmqProvisioner.ContainerExists(name)
	assert.False(exists)
}

func envVarsOk() bool {
	return testRabbitmqProv.driverConfig.DockerEndpoint != "" && testRabbitmqProv.driverConfig.DockerImage != "" && testRabbitmqProv.driverConfig.ImageVersion != ""
}
