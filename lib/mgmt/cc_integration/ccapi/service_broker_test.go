package ccapi

import (
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/pivotal-golang/lager/lagertest"
)

var loggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")
var ccApi = os.Getenv("CC_API")
var tokenEndpoint = os.Getenv("TOKEN_ENDPOINT")
var clientId = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var usbEndpoint = os.Getenv("USB_ENDPOINT")
var usbUsername = os.Getenv("USB_USERNAME")
var usbPassword = os.Getenv("USB_PASSWORD")

func TestCreate(t *testing.T) {
	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'CC_API','TOKEN_ENDPOINT','CLIENT_ID','CLIENT_SECRET','USB_ENDPOINT','USB_USERNAME','USB_PASSWORD'")
	}

	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaaapi.NewTokenGenerator(tokenEndpoint, clientId, clientSecret, client)

	sb := NewServiceBroker(client, tokenGenerator, ccApi, loggerSB)

	err := sb.Create("usb", usbEndpoint, usbUsername, usbEndpoint)
	if err != nil {
		t.Errorf("Error create service broker endpoints: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'CC_API','TOKEN_ENDPOINT','CLIENT_ID','CLIENT_SECRET','USB_ENDPOINT','USB_USERNAME','USB_PASSWORD'")
	}

	client := httpclient.NewHttpClient(true)

	tokenGenerator := uaaapi.NewTokenGenerator(tokenEndpoint, clientId, clientSecret, client)

	sb := NewServiceBroker(client, tokenGenerator, ccApi, loggerSB)

	err := sb.Update("usb", usbEndpoint, usbUsername, usbPassword)
	if err != nil {
		t.Errorf("Error update service broker endpoints: %v", err)
	}
}

func envVarsOk() bool {
	return ccApi != "" && tokenEndpoint != "" && clientId != "" && clientSecret != "" && usbEndpoint != "" && usbUsername != "" && usbPassword != ""
}
