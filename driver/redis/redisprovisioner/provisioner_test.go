package redisprovisioner

import (
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/driver/redis/config"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-provisioner")

var testRedisProv = struct {
	redisProvisioner RedisProvisionerInterface
	driverConfig     config.RedisDriverConfig
}{}

func init() {
	testRedisProv.driverConfig = config.RedisDriverConfig{
		DockerEndpoint: os.Getenv("DOCKER_ENDPOINT"),
		DockerImage:    os.Getenv("REDIS_DOCKER_IMAGE"),
		ImageVersion:   os.Getenv("REDIS_DOCKER_IMAGE_VERSION"),
	}

	testRedisProv.redisProvisioner = NewRedisProvisioner(logger)
	testRedisProv.redisProvisioner.Connect(testRedisProv.driverConfig)
}

func TestCreateContainer(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "testContainer"

	err := testRedisProv.redisProvisioner.CreateContainer(name)
	assert.NoError(err)

	exists, err := testRedisProv.redisProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.True(exists)
}

func TestGetCredentials(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "testContainer"

	credentials, err := testRedisProv.redisProvisioner.GetCredentials(name)
	assert.NoError(err)
	assert.NotNil(credentials["password"])
	assert.NotNil(credentials["port"])
}

func TestDeleteContainer(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "testContainer"

	err := testRedisProv.redisProvisioner.DeleteContainer(name)
	assert.NoError(err)

	exists, err := testRedisProv.redisProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.False(exists)
}

func envVarsOk() bool {
	return testRedisProv.driverConfig.DockerEndpoint != "" && testRedisProv.driverConfig.DockerImage != "" && testRedisProv.driverConfig.ImageVersion != ""
}
