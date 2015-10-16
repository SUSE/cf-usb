package authentication

import (
	auth "github.com/cloudfoundry-incubator/routing-api/authentication"
)

type token struct {
	accessToken  auth.Token
	uaaPublicKey string
}

func NewAccessToken(uaaPublicKey string) token {
	return token{
		accessToken: auth.NewAccessToken(uaaPublicKey),
	}
}

func (token token) DecodeToken(userToken string, desiredPermission string) error {
	return token.accessToken.DecodeToken(userToken, desiredPermission)
}

func (token token) CheckPublicToken() error {
	return token.accessToken.CheckPublicToken()
}
