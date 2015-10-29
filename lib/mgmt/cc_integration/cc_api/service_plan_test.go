package cc_api

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaa_api"
	"github.com/pivotal-golang/lager/lagertest"
)

var loggerSP *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestUpdateServicePlanVisibility(t *testing.T) {
	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaa_api.NewTokenGenerator("https://uaa.bosh-lite.com", "cc_usb_management", "cc-usb-management-secret", client)

	sp := NewServicePlan(client, tokenGenerator, "http://api.bosh-lite.com", loggerSP)

	err := sp.Update("790fb736-7be2-4aef-9770-ad51ca4b2b84", "planone", "This is the first plan", "eaa316fe-f581-44f8-99dd-ac76e3b98e42", true, true)
	if err != nil {
		t.Errorf("Error enable service access: %v", err)
	}
}
