package ccapi

import (
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var infoLogger *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestGetInfo(t *testing.T) {
	ccApi := os.Getenv("CC_API")

	if ccApi == "" {
		t.Skip("Skipping test, not all env variables are set:'CC_API'")
	}

	client := httpclient.NewHttpClient(true)

	getinfo := NewGetInfo(ccApi, client, infoLogger)

	tokenUrl, err := getinfo.GetTokenEndpoint()
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.Contains(t, tokenUrl, "uaa")
}
