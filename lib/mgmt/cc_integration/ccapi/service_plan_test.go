package ccapi

import (
	"testing"

	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var loggerSP = lagertest.NewTestLogger("cc-api")

func TestUpdateServicePlanVisibility(t *testing.T) {
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"resources":[{"metadata":{"guid":"guid"},"entity":{"name":"","free":false,"description":"","public":false,"service_guid":""}}]}`), nil)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sp := NewServicePlan(client, tokenGenerator, "http://api.1.2.3.4.io", loggerSP)

	err := sp.Update("a-service-label")
	if err != nil {
		t.Errorf("Error enable service access: %v", err)
	}

	assert.NoError(t, err)
}
