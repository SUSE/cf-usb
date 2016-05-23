package uaa

import (
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	accessToken "github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa/token"
	"github.com/pivotal-golang/lager"
)

//Auth is the structure used for storing authentication related items
type Auth struct {
	accessToken accessToken.Token
	scope       string
	logger      lager.Logger
}

//NewUaaAuth creates a new Auth with a token created from the data passed in and returns it or an error if it failes
func NewUaaAuth(uaaPublicKey string, symmetricVerificationKey string, scope string, devMode bool, logger lager.Logger) (authentication.Authentication, error) {
	var token accessToken.Token

	if devMode {
		token = accessToken.NullToken{}
	} else {
		token = accessToken.NewAccessToken(uaaPublicKey, symmetricVerificationKey)
	}

	log := logger.Session("authentication", lager.Data{"dev mode": devMode})

	newAuth := Auth{token, scope, log}

	if uaaPublicKey != "" {
		err := newAuth.accessToken.CheckPublicToken()
		if err != nil {
			return nil, err
		}
	}

	log.Debug("initializing-uaa-auth-succeeded")

	return &newAuth, nil
}

//IsAuthenticated checks if the auth header is authenticated in this scope
func (auth *Auth) IsAuthenticated(authHeader string) error {
	err := auth.accessToken.DecodeToken(authHeader, auth.scope)
	if err != nil {
		auth.logger.Error("decode-token-failed", err)
	}
	return err
}
