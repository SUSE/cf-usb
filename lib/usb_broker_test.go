package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config"
	fakes "github.com/hpcloud/cf-usb/lib/csm/fakes"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("usb-broker-test")

func setupEnv() (*UsbBroker, *fakes.FakeCSMInterface, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	buildDir := filepath.Join(workDir, "../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH), "drivers")
	os.Setenv("USB_DRIVER_PATH", buildDir)

	configFile := filepath.Join(workDir, "../test-assets/file-config/config.json")

	configProvider := config.NewFileConfig(configFile)
	fake := new(fakes.FakeCSMInterface)

	broker := NewUsbBroker(configProvider, logger, fake)
	return broker, fake, nil
}

func TestGetCatalog(t *testing.T) {
	assert := assert.New(t)
	broker, _, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}

	serviceCatalog := broker.Services()

	assert.Equal(2, len(serviceCatalog.Services))
	assert.Equal("83E94C97-C755-46A5-8653-461517EB442A", serviceCatalog.Services[0].ID)
	assert.Equal("echo", serviceCatalog.Services[0].Name)
	assert.Equal(2, len(serviceCatalog.Services[0].Plans))

	for _, plan := range serviceCatalog.Services[0].Plans {
		if plan.Name == "free" {
			assert.Equal("53425178-F731-49E7-9E53-5CF4BE9D807A", plan.ID)
			assert.Equal("This is the first plan", plan.Description)
			continue
		}
		if plan.Name == "secondary" {
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
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(false, nil)
	fake.CreateWorkspaceReturns(nil)

	provisionDetails := brokerapi.ProvisionDetails{}
	provisionDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	provisionDetails.PlanID = "53425178-F731-49E7-9E53-5CF4BE9D807A"
	provisionDetails.Parameters = make(map[string]interface{})
	provisionDetails.Parameters["param"] = "myparam"
	_, _, err = broker.Provision("newInstanceID", provisionDetails, false)
	assert.Nil(err)
}

func TestProvisionServiceNoParams(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(false, nil)
	fake.CreateWorkspaceReturns(nil)

	provisionDetails := brokerapi.ProvisionDetails{}
	provisionDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	provisionDetails.PlanID = "53425178-F731-49E7-9E53-5CF4BE9D807A"
	_, _, err = broker.Provision("newInstanceID", provisionDetails, false)
	assert.Nil(err)
}

func TestProvisionServiceInvalidParams(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(false, nil)
	fake.CreateWorkspaceReturns(errors.New("Error creating"))

	provisionDetails := brokerapi.ProvisionDetails{}
	provisionDetails.ServiceID = "83E94C97-C755-46A5-8653-461517EB442A"
	provisionDetails.PlanID = "53425178-F731-49E7-9E53-5CF4BE9D807A"
	provisionDetails.Parameters = make(map[string]interface{})
	provisionDetails.Parameters["notvalid"] = "myparam"
	_, _, err = broker.Provision("newInstanceID", provisionDetails, false)
	assert.NotNil(err)
}

func TestProvisionServiceExists(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(true, nil)

	_, _, err = broker.Provision("instanceID", brokerapi.ProvisionDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	}, false)
	assert.Equal(brokerapi.ErrInstanceAlreadyExists.Error(), err.Error())
}

func TestDeprovision(t *testing.T) {

	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(true, nil)
	fake.DeleteWorkspaceReturns(nil)
	_, err = broker.Deprovision("instanceID", brokerapi.DeprovisionDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	}, false)
	assert.Nil(err)
}

func TestDeprovisionDoesNotExist(t *testing.T) {

	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.WorkspaceExistsReturns(false, nil)

	_, err = broker.Deprovision("wrongInstanceID", brokerapi.DeprovisionDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	}, false)

	assert.NotNil(err)
	assert.Equal(brokerapi.ErrInstanceDoesNotExist.Error(), err.Error())
}

func TestBind(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	responseCreate := make(map[string]interface{})
	responseCreate["username"] = "user"
	responseCreate["password"] = "pass"
	fake.CreateConnectionReturns(responseCreate, nil)

	bindResponse, err := broker.Bind("instanceID", "newBindingID", brokerapi.BindDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	})

	response := bindResponse.Credentials.(map[string]interface{})

	assert.Equal("user", response["username"].(string))
	assert.Equal("pass", response["password"].(string))
	assert.NotNil(bindResponse)
	assert.Nil(err)
}

func TestUnbind(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.ConnectionExistsReturns(true, nil)
	fake.DeleteConnectionReturns(nil)
	err = broker.Unbind("instanceID", "credentialsID", brokerapi.UnbindDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	})

	assert.Nil(err)
}

func TestBindExists(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.ConnectionExistsReturns(true, nil)

	bindResponse, err := broker.Bind("instanceID", "credentialsID", brokerapi.BindDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	})

	assert.Nil(bindResponse.Credentials)
	assert.NotNil(err)
	assert.Equal(brokerapi.ErrBindingAlreadyExists.Error(), err.Error())
}

func TestUnbindDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	broker, fake, err := setupEnv()
	if err != nil {
		assert.Fail(err.Error())
	}
	fake.LoginReturns(nil)
	fake.ConnectionExistsReturns(false, nil)

	err = broker.Unbind("instanceID", "wrongBindingID", brokerapi.UnbindDetails{
		ServiceID: "83E94C97-C755-46A5-8653-461517EB442A",
	})

	assert.NotNil(err)
	assert.Equal(brokerapi.ErrBindingDoesNotExist.Error(), err.Error())
}
