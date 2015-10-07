package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"github.com/stretchr/testify/assert"
)

var testService = brokerapi.Service{
	ID:          "GUID",
	Name:        "testService",
	Description: "",
	Bindable:    true,
	Plans:       []brokerapi.ServicePlan{},
	Metadata:    brokerapi.ServiceMetadata{},
	Tags:        []string{},
}

var testDriverConfig = config.DriverConfig{
	DriverType: "dummy",
	ServiceIDs: []string{"GUID"},
}

var testConfig = config.Config{
	Crednetials:    brokerapi.BrokerCredentials{},
	ServiceCatalog: []brokerapi.Service{testService},
	DriverConfigs:  []config.DriverConfig{testDriverConfig},
	Listen:         ":5580",
	APIVersion:     "2.6",
	LogLevel:       "debug",
}

var testDriverProperties = config.DriverProperties{
	Services: []brokerapi.Service{testService},
}

func setupEnv() (*UsbBroker, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	buildDir := filepath.Join(wd, "../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

	configProperties := []byte(`{"property_one":"one", "property_two":"two"}`)

	testDriverConfig.Configuration = (*json.RawMessage)(&configProperties)

	testDriverProperties.DriverConfiguration = testDriverConfig.Configuration

	driverProvider, err := NewDriverProvider("dummy", testDriverProperties)
	if err != nil {
		return nil, err
	}

	broker := NewUsbBroker([]*DriverProvider{driverProvider}, &testConfig, lager.NewLogger("brokerTests"))
	return broker, nil
}

func TestGetCatalog(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	serviceCatalog := broker.Services()

	assert.Equal(1, len(serviceCatalog))
	assert.Equal("GUID", serviceCatalog[0].ID)
	assert.Equal("testService", serviceCatalog[0].Name)

	assert.Nil(err)
	assert.True(true)
}

func TestProvisionService(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Provision("newInstance", brokerapi.ProvisionDetails{
		ID: "GUID",
	})
	assert.Nil(err)
}

func TestProvisionServiceExists(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Provision("instanceID", brokerapi.ProvisionDetails{
		ID: "GUID",
	})
	assert.Equal(brokerapi.ErrInstanceAlreadyExists.Error(), err.Error())
}

func TestDeprovision(t *testing.T) {

	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Deprovision("instanceID", brokerapi.DeprovisionDetails{
		ServiceID: "GUID",
	})
	assert.Nil(err)
}

func TestDeprovisionDoesNotExist(t *testing.T) {

	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Deprovision("wrongInstanceID", brokerapi.DeprovisionDetails{
		ServiceID: "GUID",
	})

	assert.NotNil(err)
	assert.Equal(brokerapi.ErrInstanceDoesNotExist.Error(), err.Error())
}

func TestBind(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	bindResponse, err := broker.Bind("instanceID", "newBindingID", brokerapi.BindDetails{
		ServiceID: "GUID",
	})

	response := bindResponse.(map[string]interface{})

	assert.Equal("user", response["username"].(string))
	assert.Equal("pass", response["password"].(string))
	assert.NotNil(bindResponse)
	assert.Nil(err)
}

func TestUnbind(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Unbind("instanceID", "bindingID", brokerapi.UnbindDetails{
		ServiceID: "GUID",
	})

	assert.Nil(err)
}

func TestBindExists(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	bindResponse, err := broker.Bind("instanceID", "bindingID", brokerapi.BindDetails{
		ServiceID: "GUID",
	})

	assert.Nil(bindResponse)
	assert.NotNil(err)
	assert.Equal(brokerapi.ErrBindingAlreadyExists.Error(), err.Error())
}

func TestUnbindDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Unbind("instanceID", "wrongBindingID", brokerapi.UnbindDetails{
		ServiceID: "GUID",
	})

	assert.NotNil(err)
	assert.Equal(brokerapi.ErrBindingDoesNotExist.Error(), err.Error())
}
