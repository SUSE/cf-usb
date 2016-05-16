package brokertest

import (
	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/redis"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"fmt"
	"os"
	"strconv"
	"testing"
)

var orgGuid string = uuid.NewV4().String()
var spaceGuid string = uuid.NewV4().String()
var serviceGuid string = uuid.NewV4().String()
var serviceGuidAsync string = fmt.Sprintf("%[1]s-async", uuid.NewV4().String())
var serviceBindingGuid string = uuid.NewV4().String()
var instances []config.Instance

var RedisConfig = struct {
	RedisAddress  string
	RedisDatabase string
	RedisPassword string
}{}

var logger = lagertest.NewTestLogger("csm-client-test")
var csmEndpoint = ""
var authToken = ""

func init() {
	RedisConfig.RedisAddress = os.Getenv("REDIS_ADDRESS")
	RedisConfig.RedisDatabase = os.Getenv("REDIS_DATABASE")
	RedisConfig.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if RedisConfig.RedisDatabase == "" {
		RedisConfig.RedisDatabase = "0"
	}
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("csm-auth-token")
}

func setupEnv() (*lib.UsbBroker, *csm.CSMInterface, error) {
	db, err := strconv.ParseInt(RedisConfig.RedisDatabase, 10, 64)
	if err != nil {
		return nil, nil, err
	}
	provisioner, err := redis.New(RedisConfig.RedisAddress, RedisConfig.RedisPassword, db)
	if err != nil {
		return nil, nil, err
	}
	configProvider := config.NewRedisConfig(provisioner)
	csmInterface := csm.NewCSMClient(logger)
	broker := lib.NewUsbBroker(configProvider, logger, csmInterface)
	return broker, &csmInterface, nil
}

func TestBrokerAPIProvisionTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}
	if csmEndpoint == "" || authToken == "" || RedisConfig.RedisAddress == "" {
		t.Skipf("Skipping broker redis integration test - missing CSM_ENDPOINT, csm-auth-token and/or REDIS configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	details := brokerapi.ProvisionDetails{}
	response, _, err := broker.Provision(workspaceID, details, false)
	t.Log(response)
	assert.NotNil(response)
	assert.NoError(err)
}

func TestBrokerAPIBindTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}
	if csmEndpoint == "" || authToken == "" || RedisConfig.RedisAddress == "" {
		t.Skipf("Skipping broker redis integration test - missing CSM_ENDPOINT, csm-auth-token and/or REDIS configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	bindDetails := brokerapi.BindDetails{}
	response, _, err := broker.Provision(workspaceID, serviceDetails, false)
	responseBind, err := broker.Bind(workspaceID, connectionID, bindDetails)
	t.Log(response)
	assert.NotNil(response)
	assert.NotNil(responseBind)
	assert.NoError(err)
}

func TestBrokerAPIUnbindTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}
	if csmEndpoint == "" || authToken == "" || RedisConfig.RedisAddress == "" {
		t.Skipf("Skipping broker redis integration test - missing CSM_ENDPOINT, csm-auth-token and/or REDIS configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	bindDetails := brokerapi.BindDetails{}
	unbindDetails := brokerapi.UnbindDetails{}
	response, _, err := broker.Provision(workspaceID, serviceDetails, false)
	assert.NoError(err)

	responseBind, err := broker.Bind(workspaceID, connectionID, bindDetails)
	assert.NoError(err)

	err = broker.Unbind(workspaceID, connectionID, unbindDetails)
	t.Log(response)
	assert.NotNil(response)
	assert.NotNil(responseBind)
	assert.NoError(err)
}

func TestBrokerAPIDeprovisionTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}
	if csmEndpoint == "" || authToken == "" || RedisConfig.RedisAddress == "" {
		t.Skipf("Skipping broker redis integration test - missing CSM_ENDPOINT, csm-auth-token and/or REDIS configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	provisionDetails := brokerapi.ProvisionDetails{}
	deprovisionDetails := brokerapi.DeprovisionDetails{}
	response, _, err := broker.Provision(workspaceID, provisionDetails, false)
	t.Log(response)
	assert.NotNil(response)
	assert.NoError(err)

	_, err = broker.Deprovision(workspaceID, deprovisionDetails, false)
	assert.NoError(err)
}
