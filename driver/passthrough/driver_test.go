package passthrough

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("passthrough-test")
var tempFile, _ = ioutil.TempFile(os.TempDir(), "temp")
var tempPath = tempFile.Name()

func getConfig() *json.RawMessage {
	rawMessage := json.RawMessage([]byte(`{"static_config": "{\"user\": \"user\", \"password\": \"password\"}"}`))
	return &rawMessage
}

func getDriver() passthroughDriver {
	d := NewPassthroughDriver(logger)
	return d.(passthroughDriver)
}

func Test_PingFail(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	rawMessage := json.RawMessage([]byte(`{"static_config": {"user": "user"}}`))
	req := &rawMessage

	var response bool
	err := passthroughDriver.Ping(req, &response)
	assert.False(response)
	assert.Error(err)
}

func Test_Ping(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	req := getConfig()

	var response bool
	err := passthroughDriver.Ping(req, &response)
	assert.True(response)
	assert.NoError(err)
}

func Test_Provision(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()

	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.ProvisionInstanceRequest

	req.InstanceID = "testId"
	req.Config = getConfig()
	req.Dials = getConfig()

	var response usbDriver.Instance
	err := passthroughDriver.ProvisionInstance(req, &response)

	assert.NoError(err)
}

func Test_GetInstance(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.GetInstanceRequest
	req.InstanceID = "testId"
	req.Config = getConfig()
	var response usbDriver.Instance

	err := passthroughDriver.GetInstance(req, &response)

	assert.Equal(response.Status, status.Exists)
	assert.NoError(err)
}

func Test_GetDialsSchema(t *testing.T) {
	assert := assert.New(t)
	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var response string
	err := passthroughDriver.GetDailsSchema("", &response)
	assert.NoError(err)
}

func Test_GetConfigSchema(t *testing.T) {
	assert := assert.New(t)
	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var response string
	err := passthroughDriver.GetConfigSchema("", &response)
	assert.NoError(err)
}

func Test_GenerateCredentials(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.GenerateCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getConfig()
	var response interface{}

	err := passthroughDriver.GenerateCredentials(req, &response)

	assert.NoError(err)
}

func Test_GetCredentials(t *testing.T) {
	assert := assert.New(t)
	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.GetCredentialsRequest
	req.Config = getConfig()
	req.CredentialsID = "user"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := passthroughDriver.GetCredentials(req, &response)
	assert.Equal(status.Exists, response.Status)
	assert.NoError(err)
}

func Test_GetCredentialsNotFound(t *testing.T) {
	assert := assert.New(t)
	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.GetCredentialsRequest
	req.Config = getConfig()
	req.CredentialsID = "userNotFound"
	req.InstanceID = "testId"

	var response usbDriver.Credentials

	err := passthroughDriver.GetCredentials(req, &response)
	assert.Equal(status.DoesNotExist, response.Status)
	assert.NoError(err)
}

func Test_RevokeCredentials(t *testing.T) {
	assert := assert.New(t)
	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.RevokeCredentialsRequest
	req.CredentialsID = "user"
	req.InstanceID = "testId"
	req.Config = getConfig()
	var response usbDriver.Credentials

	err := passthroughDriver.RevokeCredentials(req, &response)
	assert.NoError(err)
}

func Test_Deprovision(t *testing.T) {
	assert := assert.New(t)

	passthroughDriver := getDriver()
	passthroughDriver.stateFileLocation = tempPath

	var req usbDriver.DeprovisionInstanceRequest
	req.Config = getConfig()
	req.InstanceID = "testId"

	var response usbDriver.Instance
	err := passthroughDriver.DeprovisionInstance(req, &response)

	assert.NoError(err)
}
