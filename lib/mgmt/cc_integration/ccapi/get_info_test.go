package ccapi

import (
	"testing"

	"github.com/SUSE/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var infoLogger = lagertest.NewTestLogger("cc-api")

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	client := new(mocks.HTTPClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"token_endpoint":"http://uaa.1.2.3.4.io"}`), nil)

	getinfo := NewGetInfo("http://api.1.2.3.4.io", client, infoLogger)
	tokenURL, err := getinfo.GetTokenEndpoint()
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.NoError(err)
	assert.Contains(tokenURL, "uaa")
}
