package cc_api

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaa_api"
	"github.com/pivotal-golang/lager/lagertest"
)

var loggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestCreate(t *testing.T) {
	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaa_api.NewTokenGenerator("https://uaa.bosh-lite.com", "cc_usb_management", "cc-usb-management-secret", client)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.bosh-lite.com", loggerSB)

	err := sb.Create("usb", "http://10.11.0.25:54054", "username", "password")
	if err != nil {
		t.Errorf("Error create service broker endpoints: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaa_api.NewTokenGenerator("https://uaa.bosh-lite.com", "cc_usb_management", "cc-usb-management-secret", client)

	sb := NewServiceBroker(client, tokenGenerator, "http://api.bosh-lite.com", loggerSB)

	err := sb.Update("00912776-00df-4749-a049-e80aaebb4ed9", "usb", "http://10.11.0.25:54054", "username", "password")
	if err != nil {
		t.Errorf("Error update service broker endpoints: %v", err)
	}
}
