package ccapi

import (
	"encoding/json"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var infoLogger *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestGetInfo(t *testing.T) {
	tokenEndpointMocked := GetInfoResponse{TokenEndpoint: "http://uaa.test.com"}
	values, err := json.Marshal(tokenEndpointMocked)
	if err != nil {
		t.Errorf("Error marshall token endpoint: %v", err)
	}

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(values, nil)

	getinfo := NewGetInfo("", client, infoLogger)
	tokenUrl, err := getinfo.GetTokenEndpoint()
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.NoError(t, err)
	assert.Contains(t, tokenUrl, "uaa")
}
