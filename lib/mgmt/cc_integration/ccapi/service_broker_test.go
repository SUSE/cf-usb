package ccapi

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var loggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HttpClient)
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

	client := new(mocks.HttpClient)
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

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"aguid"}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	guid, err := sb.GetServiceBrokerGuidByName("usbTest")
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

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":""},"entity":{"name":"","free":false,"description":"","public":false,"service_guid":""}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	err := sb.EnableServiceAccess("alabel")
	assert.NoError(err)
}

func TestGetServices(t *testing.T) {
	assert := assert.New(t)

	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"31f3227e-1599-44ad-bb68-1938bd4824ad","url":"/v2/services/31f3227e-1599-44ad-bb68-1938bd4824ad","created_at":"2016-01-13T17:29:41Z","updated_at":null},"entity":{"label":"label-34","provider":null,"url":null,"description":"desc-72","long_description":null,"version":null,"info_url":null,"active":true,"bindable":true,"unique_id":"a7d5938b-969b-46cf-9403-48a99fb59985","extra":null,"tags":[],"requires":[],"documentation_url":null,"service_broker_guid":"3c96db1d-962f-4301-9a2a-9fbe36b7ec50","plan_updateable":false,"service_plans_url":"/v2/services/31f3227e-1599-44ad-bb68-1938bd4824ad/service_plans"}}]}`), nil)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSB)
	assert.NotNil(sb)

	serv, err := sb.GetServices()
	t.Log(serv)
	assert.NoError(err)
}
