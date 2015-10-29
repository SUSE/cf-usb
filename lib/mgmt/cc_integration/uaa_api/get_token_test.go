package uaa_api

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
)

func TestGetToken(t *testing.T) {
	client := httpclient.NewHttpClient(true)

	tokenGenerator := NewTokenGenerator("https://uaa.bosh-lite.com", "cc_usb_management", "cc-usb-management-secret", client)

	token, err := tokenGenerator.GetToken()
	if err != nil {
		t.Errorf("Error generationg token: %v", err)
	}

	t.Logf("token: %v", token)
}
