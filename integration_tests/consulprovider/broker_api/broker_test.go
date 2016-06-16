package brokertest

import (
	"reflect"

	loads "github.com/go-openapi/loads"
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/broker"
	"github.com/hpcloud/cf-usb/lib/broker/operations"
	"github.com/hpcloud/cf-usb/lib/brokermodel"
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

func setupEnv() (*operations.BrokerAPI, error) {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if ConsulConfig.ConsulAddress == "" {
		return nil, fmt.Errorf("CONSUL configuration environment variables not set")
	}
	if csmEndpoint == "" {
		return nil, fmt.Errorf("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		return nil, fmt.Errorf("CSM_API_KEY not set")
	}

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
		return nil, err
	}
	var list api.KVPairs
	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.6")})

	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}")})

	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}")})

	err = provisioner.PutKVs(&list, nil)
	if err != nil {
		return nil, err
	}

	instanceInfo := config.Instance{}
	instanceInfo.AuthenticationKey = authToken
	instanceInfo.Name = "testInstance"
	instanceInfo.Dials = make(map[string]config.Dial)

	dialInfo := config.Dial{}
	dialInfo.Plan = brokermodel.Plan{}
	dialInfo.Plan.Free = true
	dialInfo.Plan.Name = "testPlan"
	dialID := uuid.NewV4().String()
	instanceInfo.Dials[dialID] = dialInfo

	service := brokermodel.CatalogService{}
	service.ID = "83E94C97-C755-46A5-8653-461517EB442A"
	service.Name = "testService"
	service.Plans = append(service.Plans, &dialInfo.Plan)
	instanceInfo.Service = service
	instanceInfo.TargetURL = csmEndpoint

	configProvider := config.NewConsulConfig(provisioner)
	err = configProvider.SetInstance(uuid.NewV4().String(), instanceInfo)
	if err != nil {
		return nil, err
	}

	csmInterface := csm.NewCSMClient(logger)
	swaggerSpec, err := loads.Analyzed(broker.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	brokerAPI := operations.NewBrokerAPI(swaggerSpec)

	broker.ConfigureAPI(brokerAPI, csmInterface, configProvider, logger)

	//broker := lib.NewUsbBroker(configProvider, logger, csmInterface)
	return brokerAPI, nil
}

func TestBrokerAPIProvisionTest(t *testing.T) {
	assert := assert.New(t)
	brokerA, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}

	workspaceID := uuid.NewV4().String()
	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	params.InstanceID = workspaceID

	response := brokerA.CreateServiceInstanceHandler.Handle(params, false)
	assert.NotNil(response)
	if reflect.TypeOf(response).String() == "*operations.CreateServiceInstanceDefault" {
		resp := response.(*operations.CreateServiceInstanceDefault)
		assert.Fail(*resp.Payload.Message)
		return
	}

	assert.Equal(
		reflect.TypeOf(operations.CreateServiceInstanceCreated{}),
		reflect.ValueOf(response).Elem().Type(),
		"Wrong response type while binding")
}

func TestBrokerAPIBindTest(t *testing.T) {
	assert := assert.New(t)
	brokerA, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}

	workspaceID := uuid.NewV4().String()
	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	params.InstanceID = workspaceID
	response := brokerA.CreateServiceInstanceHandler.Handle(params, false)

	connectionID := uuid.NewV4().String()
	if assert.NotNil(response) {
		assert.Equal("*operations.CreateServiceInstanceCreated", reflect.TypeOf(response).String())
		connParams := operations.ServiceBindParams{}
		connParams.InstanceID = workspaceID
		connParams.BindingID = connectionID
		connParams.Binding = &brokermodel.Binding{}
		connParams.Binding.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
		resp := brokerA.ServiceBindHandler.Handle(connParams, false)
		assert.NotNil(resp)
		t.Log(reflect.ValueOf(resp).Elem().Type())

		switch resp.(type) {
		case *operations.ServiceBindCreated:
			break
		case *operations.ServiceBindDefault:
			assert.FailNow("Waiting for ServiceBindCreated, but got ServiceBindDefault")
			resp := response.(*operations.ServiceBindDefault)
			assert.Fail(*resp.Payload.Message)
			return
		default:
			assert.Fail("No error response should happen")
			return
		}

		assert.Equal(reflect.TypeOf(operations.ServiceBindCreated{}),
			reflect.ValueOf(resp).Elem().Type())
	}
}

func TestBrokerAPIUnbindTest(t *testing.T) {
	assert := assert.New(t)
	brokerA, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}

	workspaceID := uuid.NewV4().String()
	connectionID := uuid.NewV4().String()

	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	params.InstanceID = workspaceID

	response := brokerA.CreateServiceInstanceHandler.Handle(params, false)
	if assert.NotNil(response) &&
		assert.Equal(
			reflect.TypeOf(operations.CreateServiceInstanceCreated{}),
			reflect.ValueOf(response).Elem().Type(),
			"Wrong response type while binding") {
		assert.Equal("*operations.CreateServiceInstanceCreated", reflect.TypeOf(response).String())
		connParams := operations.ServiceBindParams{}
		connParams.InstanceID = workspaceID
		connParams.BindingID = connectionID
		connParams.Binding = &brokermodel.Binding{}
		connParams.Binding.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
		resp := brokerA.ServiceBindHandler.Handle(connParams, false)
		if assert.NotNil(resp, "There should be an answer when binding") && assert.Equal(reflect.TypeOf(operations.ServiceBindCreated{}), reflect.ValueOf(resp).Elem().Type(), "Wrong response type while binding") {
			unbindParams := operations.ServiceUnbindParams{}
			unbindParams.InstanceID = workspaceID
			unbindParams.BindingID = connectionID
			unbindParams.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
			respUnbind := brokerA.ServiceUnbindHandler.Handle(unbindParams, false)
			if assert.NotNil(respUnbind, "There should be an unswer when unbinding") {
				switch respUnbind.(type) {
				case *operations.ServiceUnbindOK:
					break
				case *operations.ServiceUnbindDefault:
					assert.Fail("Waiting for ServiceUnbindOK, but Got ServiceUnbindDefault")
					resp := response.(*operations.ServiceUnbindDefault)
					assert.Fail(*resp.Payload.Message)
					break
				case *operations.ServiceUnbindGone:
					assert.Fail("Waiting for ServiceUnbindOK, but Got ServiceUnbindGone")
					break
				default:
					assert.FailNow("No error response should happen")
				}
			}
		}

	}

}

func TestBrokerAPIDeprovisionTest(t *testing.T) {
	assert := assert.New(t)
	brokerA, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}

	workspaceID := uuid.NewV4().String()
	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	params.InstanceID = workspaceID
	deprovisionParams := operations.DeprovisionServiceInstanceParams{}
	deprovisionParams.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	deprovisionParams.InstanceID = workspaceID

	response := brokerA.CreateServiceInstanceHandler.Handle(params, false)
	t.Log(response)
	if assert.NotNil(response, "There should be an answer when provisioning") &&
		assert.Equal(
			reflect.TypeOf(operations.CreateServiceInstanceCreated{}),
			reflect.ValueOf(response).Elem().Type(),
			"Wrong response type while binding") {
		resp := brokerA.DeprovisionServiceInstanceHandler.Handle(deprovisionParams, false)
		if assert.NotNil(resp, "There should be an unswer when unprovisioning") {
			switch resp.(type) {
			case *operations.DeprovisionServiceInstanceOK:
				break
			case *operations.DeprovisionServiceInstanceDefault:
				assert.Fail("Waiting for DeprovisionServiceInstanceOK, but Got DeprovisionServiceInstanceDefault")
				break
			case *operations.DeprovisionServiceInstanceGone:
				assert.Fail("Waiting for DeprovisionServiceInstanceOK, but Got DeprovisionServiceInstanceGone")
				break
			default:
				assert.FailNow("No error response should happen when deprovisioning")
			}
		}
	}
}
