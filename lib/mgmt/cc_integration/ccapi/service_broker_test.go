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
