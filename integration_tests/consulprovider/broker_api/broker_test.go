package brokertest

import (
	"github.com/frodenas/brokerapi"
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"fmt"
	"os"
	"testing"
)

var orgGuid string = uuid.NewV4().String()
var spaceGuid string = uuid.NewV4().String()
var serviceGuid string = uuid.NewV4().String()
var serviceGuidAsync string = fmt.Sprintf("%[1]s-async", uuid.NewV4().String())
var serviceBindingGuid string = uuid.NewV4().String()
var instances []config.Instance

var ConsulConfig = struct {
	ConsulAddress    string
	ConsulDatacenter string
	ConsulUser       string
	ConsulPassword   string
	ConsulSchema     string
	ConsulToken      string
}{}

var logger = lagertest.NewTestLogger("csm-client-test")
var csmEndpoint = ""
var authToken = ""

func init() {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
}

func setupEnv() (*lib.UsbBroker, *csm.Interface, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.ConsulAddress
	consulConfig.Datacenter = ConsulConfig.ConsulDatacenter

	var auth api.HttpBasicAuth
	auth.Username = ConsulConfig.ConsulUser
	auth.Password = ConsulConfig.ConsulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = ConsulConfig.ConsulSchema

	consulConfig.Token = ConsulConfig.ConsulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return nil, nil, err
	}
	var list api.KVPairs
	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.6")})

	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}")})

	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}")})

	err = provisioner.PutKVs(&list, nil)
	if err != nil {
		return nil, nil, err
	}

	instanceInfo := config.Instance{}
	instanceInfo.AuthenticationKey = authToken
	instanceInfo.Name = "testInstance"
	instanceInfo.Dials = make(map[string]config.Dial)

	dialInfo := config.Dial{}
	dialInfo.Plan = brokerapi.ServicePlan{}
	dialInfo.Plan.Free = true
	dialInfo.Plan.Name = "testPlan"
	dialID := uuid.NewV4().String()
	instanceInfo.Dials[dialID] = dialInfo

	service := brokerapi.Service{}
	service.ID = "83E94C97-C755-46A5-8653-461517EB442A"
	service.Name = "testService"
	service.Plans = append(service.Plans, dialInfo.Plan)
	instanceInfo.Service = service
	instanceInfo.TargetURL = csmEndpoint

	configProvider := config.NewConsulConfig(provisioner)
	err = configProvider.SetInstance(uuid.NewV4().String(), instanceInfo)
	if err != nil {
		return nil, nil, err
	}

	csmInterface := csm.NewCSMClient(logger)
	broker := lib.NewUsbBroker(configProvider, logger, csmInterface)
	return broker, &csmInterface, nil
}

func TestBrokerAPIProvisionTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Error(err)
	}
	if csmEndpoint == "" || authToken == "" || ConsulConfig.ConsulAddress == "" {
		t.Skipf("Skipping broker consul integration test - missing CSM_ENDPOINT, CSM_API_KEY and/or CONSUL configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	details := brokerapi.ProvisionDetails{}
	details.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	response, _, err := broker.Provision(workspaceID, details, false)
	t.Log(response)
	assert.NotNil(response)
	assert.NoError(err)
}

func TestBrokerAPIBindTest(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		t.Error(err)
	}
	if csmEndpoint == "" || authToken == "" || ConsulConfig.ConsulAddress == "" {
		t.Skipf("Skipping broker consul integration test - missing CSM_ENDPOINT, CSM_API_KEY and/or CONSUL configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	serviceDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	bindDetails := brokerapi.BindDetails{}
	bindDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
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
		t.Error(err)
	}
	if csmEndpoint == "" || authToken == "" || ConsulConfig.ConsulAddress == "" {
		t.Skipf("Skipping broker consul integration test - missing CSM_ENDPOINT, CSM_API_KEY and/or CONSUL configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()
	serviceDetails := brokerapi.ProvisionDetails{}
	serviceDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	bindDetails := brokerapi.BindDetails{}
	bindDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	unbindDetails := brokerapi.UnbindDetails{}
	unbindDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
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
		t.Error(err)
	}
	if csmEndpoint == "" || authToken == "" || ConsulConfig.ConsulAddress == "" {
		t.Skipf("Skipping broker consul integration test - missing CSM_ENDPOINT, CSM_API_KEY and/or CONSUL configuration environment variables")
	}

	workspaceID := uuid.NewV4().String()
	provisionDetails := brokerapi.ProvisionDetails{}
	provisionDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	deprovisionDetails := brokerapi.DeprovisionDetails{}
	deprovisionDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	response, _, err := broker.Provision(workspaceID, provisionDetails, false)
	t.Log(response)
	assert.NotNil(response)
	assert.NoError(err)

	_, err = broker.Deprovision(workspaceID, deprovisionDetails, false)
	assert.NoError(err)
}
