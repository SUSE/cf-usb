package cc_api

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var infoLogger *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestGetInfo(t *testing.T) {
	client := httpclient.NewHttpClient(true)

	getinfo := NewGetInfo("http://api.bosh-lite.com", client, infoLogger)

	tokenUrl, err := getinfo.GetTokenEndpoint()
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.Contains(t, tokenUrl, "uaa")
}
