package lib

import (
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

func setupEnv() (*UsbBroker, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	buildDir := filepath.Join(workDir, "../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

	var logger = lager.NewLogger("test")

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

	configFile := filepath.Join(workDir, "../test-assets/file-config/dummy_config.json")

	configProvider := config.NewFileConfig(configFile)

	testDriverProperties, err := configProvider.GetDriverInstanceConfig("A0000000-0000-0000-0000-000000000002")
	if err != nil {
		return nil, err
	}

	driverProvider, err := NewDriverProvider("dummy", *testDriverProperties, logger)
	if err != nil {
		return nil, err
	}

	testConfig, err := configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	broker := NewUsbBroker([]*DriverProvider{driverProvider},
		testConfig, lager.NewLogger("brokerTests"))
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
	assert.Equal(2, len(serviceCatalog[0].Plans))

	for _, plan := range serviceCatalog[0].Plans {
		if plan.Name == "planone" {
			assert.Equal("53425178-F731-49E7-9E53-5CF4BE9D807A", plan.ID)
			assert.Equal("This is the first plan", plan.Description)
			continue
		}
		if plan.Name == "plantwo" {
			assert.Equal("888B59E0-C2A1-4AB6-9335-2E90114A8F07", plan.ID)
			assert.Equal("This is the secondary plan", plan.Description)
			continue
		}
		assert.Fail("Plans are not parsed correctly")
	}

	assert.Nil(err)
	assert.True(true)
}

func TestProvisionService(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = broker.Provision("newInstanceID", brokerapi.ProvisionDetails{
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

	err = broker.Unbind("instanceID", "credentialsID", brokerapi.UnbindDetails{
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

	bindResponse, err := broker.Bind("instanceID", "credentialsID", brokerapi.BindDetails{
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
