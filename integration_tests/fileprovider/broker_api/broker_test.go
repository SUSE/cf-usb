package brokertest

import (
	"reflect"

	loads "github.com/go-openapi/loads"
	"github.com/hpcloud/cf-usb/lib/broker"
	"github.com/hpcloud/cf-usb/lib/broker/operations"
	"github.com/hpcloud/cf-usb/lib/brokermodel"
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

var logger = lagertest.NewTestLogger("csm-client-test")
var csmEndpoint = ""
var authToken = ""
var serviceID = ""

func setupEnv() (*operations.BrokerAPI, error) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" {
		return nil, fmt.Errorf("CSM_ENDPOINT not set")
	}
	if authToken == "" {
		return nil, fmt.Errorf("CSM_API_KEY not set")
	}
	file, err := ioutil.TempFile(os.TempDir(), "brokertest")
	if err != nil {
		return nil, err
	}
	workDir, err := os.Getwd()
	testFile := filepath.Join(workDir, "../../../test-assets/file-config/config.json")

	info, err := ioutil.ReadFile(testFile)
	if err != nil {
		return nil, err
	}
	content := string(info)
	content = strings.Replace(content, "http://127.0.0.1:8080", csmEndpoint, -1)
	content = strings.Replace(content, "authkey", authToken, -1)

	configFile := file.Name()

	err = ioutil.WriteFile(configFile, []byte(content), 0777)
	if err != nil {
		return nil, err
	}

	configProvider := config.NewFileConfig(configFile)
	config, err := configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, instance := range config.Instances {
		serviceID = instance.Service.ID
		break
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}

	workspaceID := uuid.NewV4().String()

	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = serviceID
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}

	workspaceID := uuid.NewV4().String()
	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = serviceID
	params.InstanceID = workspaceID
	response := brokerA.CreateServiceInstanceHandler.Handle(params, false)

	connectionID := uuid.NewV4().String()
	if assert.NotNil(response) {
		assert.Equal("*operations.CreateServiceInstanceCreated", reflect.TypeOf(response).String())
		connParams := operations.ServiceBindParams{}
		connParams.InstanceID = workspaceID
		connParams.BindingID = connectionID
		connParams.Binding = &brokermodel.Binding{}
		connParams.Binding.ServiceID = serviceID
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

		assert.Equal(reflect.TypeOf(operations.ServiceBindCreated{}), reflect.ValueOf(resp).Elem().Type()) //reflect.TypeOf(resp).String())
	}
}

func TestBrokerAPIUnbindTest(t *testing.T) {
	assert := assert.New(t)
	brokerA, err := setupEnv()
	if err != nil {
		t.Skip(err)
	}
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
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
		connParams.Binding.ServiceID = serviceID
		resp := brokerA.ServiceBindHandler.Handle(connParams, false)
		if assert.NotNil(resp, "There should be an answer when binding") && assert.Equal(reflect.TypeOf(operations.ServiceBindCreated{}), reflect.ValueOf(resp).Elem().Type(), "Wrong response type while binding") {
			unbindParams := operations.ServiceUnbindParams{}
			unbindParams.InstanceID = workspaceID
			unbindParams.BindingID = connectionID
			unbindParams.ServiceID = serviceID
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
	if serviceID == "" {
		t.Skip("Config file does not contain a service definition")
	}
	workspaceID := uuid.NewV4().String()
	params := operations.CreateServiceInstanceParams{}
	params.Service = &brokermodel.Service{}
	params.Service.ServiceID = serviceID
	params.InstanceID = workspaceID
	deprovisionParams := operations.DeprovisionServiceInstanceParams{}
	deprovisionParams.ServiceID = serviceID
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
