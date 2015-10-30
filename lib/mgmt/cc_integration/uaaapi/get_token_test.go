package uaaapi

import (
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
)

var tokenEndpoint = os.Getenv("TOKEN_ENDPOINT")
var clientId = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

func TestGetToken(t *testing.T) {
	if !envVarsOk() {
		t.Skip("Skipping test, not all env variables are set:'TOKEN_ENDPOINT','CLIENT_ID','CLIENT_SECRET'")
	}

	client := httpclient.NewHttpClient(true)

	tokenGenerator := NewTokenGenerator(tokenEndpoint, clientId, clientSecret, client)

	token, err := tokenGenerator.GetToken()
	if err != nil {
		t.Errorf("Error generationg token: %v", err)
	}

	t.Logf("token: %v", token)
}

func envVarsOk() bool {
	return tokenEndpoint != "" && clientId != "" && clientSecret != ""
}
