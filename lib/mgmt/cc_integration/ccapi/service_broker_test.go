package ccapi

import (
	"encoding/json"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var loggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestCreate(t *testing.T) {
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sb := NewServiceBroker(client, tokenGenerator, "ccApi", loggerSB)

	err := sb.Create("usb", "brokerUrl", "brokerUsername", "brokerPassword")
	if err != nil {
		t.Errorf("Error create service broker endpoints: %v", err)
	}

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	var bra []BrokerResource
	br := BrokerResource{Values: BrokerMetadata{Guid: "aguid"}}
	bra = append(bra, br)

	getBrokerMocked := BrokerResources{Resources: bra}
	values, err := json.Marshal(getBrokerMocked)
	if err != nil {
		t.Errorf("Error marshall get broker: %v", err)
	}

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(values, nil)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sb := NewServiceBroker(client, tokenGenerator, "ccApi", loggerSB)
	err = sb.Update("usb", "brokerUrl", "brokerUsername", "brokerPassword")
	if err != nil {
		t.Errorf("Error update service broker endpoints: %v", err)
	}

	assert.NoError(t, err)
}
