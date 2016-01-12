package uaa

import (
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	accessToken "github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa/token"
	"github.com/pivotal-golang/lager"
)

type UaaAuth struct {
	accessToken accessToken.Token
	scope       string
	logger      lager.Logger
}

func NewUaaAuth(uaaPublicKey string, symmetricVerificationKey string, scope string, devMode bool, logger lager.Logger) (authentication.AuthenticationInterface, error) {
	var token accessToken.Token

	if devMode {
		token = accessToken.NullToken{}
	} else {
		token = accessToken.NewAccessToken(uaaPublicKey, symmetricVerificationKey)
	}

	log := logger.Session("authentication", lager.Data{"dev mode": devMode})

	newAuth := UaaAuth{token, scope, log}

	if uaaPublicKey != "" {
		err := newAuth.accessToken.CheckPublicToken()
		if err != nil {
			return nil, err
		}
	}

	log.Debug("initializing-uaa-auth-succeeded")

	return &newAuth, nil
}

func (auth *UaaAuth) IsAuthenticated(authHeader string) error {
	err := auth.accessToken.DecodeToken(authHeader, auth.scope)
	if err != nil {
		auth.logger.Error("decode-token-failed", err)
	}
	return err
}
