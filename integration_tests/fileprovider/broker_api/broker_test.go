package brokertest

import (
	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var orgGuid string = uuid.NewV4().String()
var spaceGuid string = uuid.NewV4().String()
var serviceGuid string = uuid.NewV4().String()
var serviceGuidAsync string = fmt.Sprintf("%[1]s-async", uuid.NewV4().String())
var serviceBindingGuid string = uuid.NewV4().String()
var instances []config.Instance

var logger = lagertest.NewTestLogger("csm-client-test")
var csmEndpoint = ""
var authToken = ""
var serviceID = ""

func setupEnv() (*lib.UsbBroker, *csm.CSM, error) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" {
		return nil, nil, fmt.Errorf("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		return nil, nil, fmt.Errorf("CSM_API_KEY not set")
	}
	file, err := ioutil.TempFile(os.TempDir(), "brokertest")
	if err != nil {
		return nil, nil, err
	}
	workDir, err := os.Getwd()
	testFile := filepath.Join(workDir, "../../../test-assets/file-config/config.json")

	info, err := ioutil.ReadFile(testFile)
	if err != nil {
		return nil, nil, err
	}
	content := string(info)
	content = strings.Replace(content, "http://127.0.0.1:8080", csmEndpoint, -1)
	content = strings.Replace(content, "authkey", authToken, -1)

	configFile := file.Name()

	err = ioutil.WriteFile(configFile, []byte(content), 0777)
	if err != nil {
		return nil, nil, err
	}

	configProvider := config.NewFileConfig(configFile)
	config, err := configProvider.LoadConfiguration()
	if err != nil {
		return nil, nil, err
	}
	for _, instance := range config.Instances {
		serviceID = instance.Service.ID
		break
	}
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}

	workspaceID := uuid.NewV4().String()
	details := brokerapi.ProvisionDetails{}
	details.ServiceID = serviceID
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}
	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	serviceDetails.ServiceID = serviceID
	bindDetails := brokerapi.BindDetails{}
	bindDetails.ServiceID = serviceID
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	serviceDetails.ServiceID = serviceID
	bindDetails := brokerapi.BindDetails{}
	bindDetails.ServiceID = serviceID
	unbindDetails := brokerapi.UnbindDetails{}
	unbindDetails.ServiceID = serviceID
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}
	workspaceID := uuid.NewV4().String()
	provisionDetails := brokerapi.ProvisionDetails{}
	provisionDetails.ServiceID = serviceID
	deprovisionDetails := brokerapi.DeprovisionDetails{}
	deprovisionDetails.ServiceID = serviceID
	response, _, err := broker.Provision(workspaceID, provisionDetails, false)
	t.Log(response)
	assert.NotNil(response)
	assert.NoError(err)

	_, err = broker.Deprovision(workspaceID, deprovisionDetails, false)
	assert.NoError(err)
}
