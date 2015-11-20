package uaa

import (
	auth "github.com/cloudfoundry-incubator/routing-api/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/pivotal-golang/lager"
)

type UaaAuth struct {
	accessToken auth.Token
	scope       string
	logger      lager.Logger
}

func NewUaaAuth(uaaPublicKey string, scope string, devMode bool, logger lager.Logger) (authentication.AuthenticationInterface, error) {
	var token auth.Token

	if devMode {
		token = auth.NullToken{}
	} else {
		token = auth.NewAccessToken(uaaPublicKey)
	}

	newAuth := UaaAuth{token, scope, logger}
	err := newAuth.accessToken.CheckPublicToken()
	if err != nil {
		return nil, err
	}
	return &newAuth, nil
}

func (auth *UaaAuth) IsAuthenticated(authHeader string) error {
	err := auth.accessToken.DecodeToken(authHeader, auth.scope)
	if err != nil {
		auth.logger.Error("is-authenticated", err)
	}
	return err
}
