package ccapi

import (
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager/lagertest"
)

var loggerSP *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestUpdateServicePlanVisibility(t *testing.T) {
	ccApi = os.Getenv("CC_API")
	tokenEndpoint = os.Getenv("TOKEN_ENDPOINT")
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	if ccApi == "" || tokenEndpoint == "" || clientId == "" || clientSecret == "" {
		t.Skip("Skipping test, not all env variables are set:'CC_API','TOKEN_ENDPOINT','CLIENT_ID','CLIENT_SECRET'")
	}

	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaaapi.NewTokenGenerator(tokenEndpoint, clientId, clientSecret, client)

	sp := NewServicePlan(client, tokenGenerator, ccApi, loggerSP)

	err := sp.Update("eaa316fe-f581-44f8-99dd-ac76e3b98e42")
	if err != nil {
		t.Errorf("Error enable service access: %v", err)
	}
}
