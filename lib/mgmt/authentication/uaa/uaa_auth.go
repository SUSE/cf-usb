package uaa

import (
	auth "github.com/cloudfoundry-incubator/routing-api/authentication"

	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
)

type UaaAuth struct {
	accessToken auth.Token
	scope       string
}

func NewUaaAuth(uaaPublicKey string, scope string, devMode bool) (authentication.AuthenticationInterface, error) {
	var token auth.Token

	if devMode {
		token = auth.NullToken{}
	} else {
		token = auth.NewAccessToken(uaaPublicKey)
	}

	newAuth := UaaAuth{token, scope}
	err := newAuth.accessToken.CheckPublicToken()
	if err != nil {
		return nil, err
	}
	return &newAuth, nil
}

func (auth *UaaAuth) IsAuthenticated(authHeader string) error {
	return auth.accessToken.DecodeToken(authHeader, auth.scope)
}
