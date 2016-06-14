package brokertest

import (
	"reflect"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/broker"
	"github.com/hpcloud/cf-usb/lib/broker/operations"
	"github.com/hpcloud/cf-usb/lib/brokermodel"
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
	"time"
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

func setupEnv() (*operations.BrokerAPI, error) {
	RedisConfig.RedisAddress = os.Getenv("REDIS_ADDRESS")
	RedisConfig.RedisDatabase = os.Getenv("REDIS_DATABASE")
	RedisConfig.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if RedisConfig.RedisDatabase == "" {
		RedisConfig.RedisDatabase = "0"
	}
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if RedisConfig.RedisAddress == "" {
		return nil, fmt.Errorf("REDIS configuration environment variables not set")
	}
	if csmEndpoint == "" {
		return nil, fmt.Errorf("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		return nil, fmt.Errorf("CSM_API_KEY not set")
	}
	db, err := strconv.ParseInt(RedisConfig.RedisDatabase, 10, 64)
	if err != nil {
		return nil, err
	}
	provisioner, err := redis.New(RedisConfig.RedisAddress, RedisConfig.RedisPassword, db)
	if err != nil {
		return nil, err
	}

	err = provisioner.SetKV("usb", "{\"api_version\":\"2.6\",\"logLevel\":\"debug\",\"broker_api\":{\"external_url\":\"http://1.2.3.4:54054\",\"listen\":\":54054\",\"credentials\":{\"username\":\"username\",\"password\":\"password\"}},\"routes_register\":{\"nats_members\":[\"nats1\",\"nats2\"],\"broker_api_host\":\"broker\",\"management_api_host\":\"management\"},\"management_api\":{\"listen\":\":54053\",\"dev_mode\":false,\"broker_name\":\"usb\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"-----BEGIN PUBLIC KEY-----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB-----END PUBLIC KEY-----\"}},\"cloud_controller\":{\"api\":\"http://api.bosh-lite.com\",\"skip_tls_validation\":true}},\"instances\":{\"A0000000-0000-0000-0000-000000000002\":{\"name\":\"dummy1\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000001\":{\"configuration\":{\"max_dbsize_mb\":2},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807A\",\"description\":\"This is the first plan\",\"free\":true}},\"B0000000-0000-0000-0000-000000000002\":{\"configuration\":{\"max_dbsize_mb\":100},\"plan\":{\"name\":\"secondary\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F07\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461512EB442A\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}},\"A0000000-0000-0000-0000-000000000003\":{\"name\":\"dummy2\",\"target\":\"http://127.0.0.1:8080\",\"authentication_key\":\"authkey\",\"dials\":{\"B0000000-0000-0000-0000-000000000011\":{\"plan\":{\"name\":\"plandummy2\",\"id\":\"888B59E0-C2A1-4AB6-9335-2E90114A8F01\",\"description\":\"This is the secondary plan\",\"free\":false}}},\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"},\"service\":{\"id\":\"83E94C97-C755-46A5-8653-461517EB442B\",\"bindable\":true,\"name\":\"echo\",\"description\":\"echo Service\",\"tags\":[\"echo\"],\"metadata\":{\"providerDisplayName\":\"Echo Service Ltd.\"}}}}}", 5*time.Minute)
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

	configProvider := config.NewRedisConfig(provisioner)
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
