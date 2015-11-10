package ccapi

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var infoLogger *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"token_endpoint":"http://uaa.1.2.3.4.io"}`), nil)

	getinfo := NewGetInfo("http://api.1.2.3.4.io", client, infoLogger)
	tokenUrl, err := getinfo.GetTokenEndpoint()
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.NoError(err)
	assert.Contains(tokenUrl, "uaa")
}
