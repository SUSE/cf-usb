package ccapi

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var loggerSB = lagertest.NewTestLogger("cc-api")

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	err := sb.Create("usbTest", "http://1.2.3.4:54054", "brokerUsername", "brokerPassword")
	if err != nil {
		t.Errorf("Error create service broker: %v", err)
	}

	assert.NoError(err)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.bosh-lite.com", loggerSB)
	assert.NotNil(sb)

	err := sb.Update("a-broker-guid", "usbTest", "http://1.2.3.4:54054", "brokerUsername", "brokerPassword")
	if err != nil {
		t.Errorf("Error update service broker: %v", err)
	}

	assert.NoError(err)
}

func TestGetServiceBrokerGuidByName(t *testing.T) {
	assert := assert.New(t)
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"aguid"}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	guid, err := sb.GetServiceBrokerGUIDByName("usbTest")
	if err != nil {
		t.Errorf("Error get service broker by name: %v", err)
	}

	assert.NoError(err)
	assert.Equal("aguid", guid)
}

func TestEnableServiceAccess(t *testing.T) {
	assert := assert.New(t)

	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"guid"},"entity":{"name":"","free":false,"description":"","public":false,"service_guid":""}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	err := sb.EnableServiceAccess("alabel")
	assert.NoError(err)
}

func TestCheckServiceNameExists(t *testing.T) {
	assert := assert.New(t)

	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"688f14a3-a5fc-4fa3-bc82-07338c180f64","url":"/v2/services/688f14a3-a5fc-4fa3-bc82-07338c180f64","created_at":"2016-01-19T19:41:27Z","updated_at":null},"entity":{"label":"label-66","provider":null,"url":null,"description":"desc-214","long_description":null,"version":null,"info_url":null,"active":true,"bindable":true,"unique_id":"641552b1-b35b-402f-81f0-6c27a3677427","extra":null,"tags":[],"requires":[],"documentation_url":null,"service_broker_guid":"aa8c13b3-f362-4d42-b25a-1be9a519e425","plan_updateable":false,"service_plans_url":"/v2/services/688f14a3-a5fc-4fa3-bc82-07338c180f64/service_plans"}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	exists, err := sb.CheckServiceNameExists("label-66")
	assert.NoError(err)
	assert.True(exists)
}
