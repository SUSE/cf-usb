package uaaapi

import (
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetToken(t *testing.T) {
	assert := assert.New(t)
	tokenValue := "atoken"

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"access_token":"atoken","expires_in":10000}`), nil)

	tokenGenerator := NewTokenGenerator("http://api.1.2.3.4.io", "clientId", "clientSecret", client)

	token, err := tokenGenerator.GetToken()
	if err != nil {
		t.Errorf("Error generationg token: %v", err)
	}

	assert.NoError(err)
	assert.Equal("bearer "+tokenValue, token)
}

func TestGetWrongToken(t *testing.T) {
	assert := assert.New(t)

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return([]byte(`{"access_token":"atoken","expires_in":""}`), nil)

	tokenGenerator := NewTokenGenerator("http://api.1.2.3.4.io", "clientId", "clientSecret", client)

	token, err := tokenGenerator.GetToken()

	assert.Error(err, "json: cannot unmarshal string into Go value of type int")
	assert.Equal("", token)
}
