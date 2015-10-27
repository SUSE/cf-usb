package redisprovisioner

import (
	"os"
	"testing"

	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-provisioner")

var testRedisProv = struct {
	redisProvisioner       RedisProvisionerInterface
	redisServiceProperties RedisServiceProperties
}{}

func init() {
	testRedisProv.redisServiceProperties = RedisServiceProperties{
		DockerEndpoint: os.Getenv("DOCKER_ENDPOINT"),
		DockerImage:    os.Getenv("DOCKER_IMAGE"),
		ImageVersion:   os.Getenv("DOCKER_IMAGE_VERSION"),
	}

	testRedisProv.redisProvisioner = NewRedisProvisioner(testRedisProv.redisServiceProperties, logger)
	testRedisProv.redisProvisioner.Init()
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
	return testRedisProv.redisServiceProperties.DockerEndpoint != "" && testRedisProv.redisServiceProperties.DockerImage != "" && testRedisProv.redisServiceProperties.ImageVersion != ""
}