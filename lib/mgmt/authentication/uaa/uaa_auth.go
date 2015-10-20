package uaaauthentication

import (
	auth "github.com/cloudfoundry-incubator/routing-api/authentication"

	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
)

type UaaAuth struct {
	accessToken auth.Token
	scope       string
}

func NewUaaAuth(uaaPublicKey string, scope string) (authentication.AuthenticationInterface, error) {
	auth := UaaAuth{auth.NewAccessToken(uaaPublicKey), scope}
	err := auth.accessToken.CheckPublicToken()
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (auth *UaaAuth) IsAuthenticated(authHeader string) error {
	return auth.accessToken.DecodeToken(authHeader, auth.scope)
}
