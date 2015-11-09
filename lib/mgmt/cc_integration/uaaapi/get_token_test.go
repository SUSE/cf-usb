package uaaapi

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var tokenEndpoint = os.Getenv("TOKEN_ENDPOINT")
var clientId = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

func TestGetToken(t *testing.T) {
	tokenValue := "atoken"
	tokenMocked := Token{AccessToken: tokenValue, ExpireTime: 10000}
	values, err := json.Marshal(tokenMocked)
	if err != nil {
		t.Errorf("Error marshall token: %v", err)
	}

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(values, nil)

	tokenGenerator := NewTokenGenerator("", "", "", client)

	token, err := tokenGenerator.GetToken()
	if err != nil {
		t.Errorf("Error generationg token: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, "bearer "+tokenValue, token)
}
