package token

// Package forked from here:
// https://github.com/cloudfoundry-incubator/routing-api/tree/877339530a78bfd01a8009fc689bca3b327a3d77/authentication

import (
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//Token is used for UAA authentication
type Token interface {
	DecodeToken(userToken string, desiredPermissions ...string) error
	CheckPublicToken() error
}

//NullToken is a token used in develop mode
type NullToken struct{}

//DecodeToken -
func (NT NullToken) DecodeToken(r string, r1 ...string) error {
	return nil
}

//CheckPublicToken for a null token the public token will never return an error
func (NT NullToken) CheckPublicToken() error {
	return nil
}

//AccessToken is the definition of an Access Token
type AccessToken struct {
	uaaPublicKey                string
	uaaSymmetricVerificationKey string
}

//NewAccessToken creates a new access token
func NewAccessToken(uaaPublicKey string, uaaSymmetricVerificationKey string) AccessToken {
	return AccessToken{
		uaaPublicKey:                uaaPublicKey,
		uaaSymmetricVerificationKey: uaaSymmetricVerificationKey,
	}
}

//DecodeToken checks if a userToken has the desired permissionss
func (accessToken AccessToken) DecodeToken(userToken string, desiredPermissions ...string) error {
	userToken, err := checkTokenFormat(userToken)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(userToken, func(t *jwt.Token) (interface{}, error) {
		if accessToken.uaaPublicKey != "" {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); ok {
				return []byte(accessToken.uaaPublicKey), nil
			}
		}

		if accessToken.uaaSymmetricVerificationKey != "" {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
				return []byte(accessToken.uaaSymmetricVerificationKey), nil
			}
		}

		return nil, fmt.Errorf("Unsupported signing method: %v", t.Header["alg"])
	})

	if err != nil {
		return err
	}

	hasPermission := false
	permissions := token.Claims["scope"]

	a := permissions.([]interface{})

	for _, permission := range a {
		for _, desiredPermission := range desiredPermissions {
			if permission.(string) == desiredPermission {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		err = errors.New("Token does not have '" + strings.Join(desiredPermissions, "', '") + "' scope")
		return err
	}

	return nil
}

//CheckPublicToken checks the validity of the public token
func (accessToken AccessToken) CheckPublicToken() error {
	var block *pem.Block
	if block, _ = pem.Decode([]byte(accessToken.uaaPublicKey)); block == nil {
		return errors.New("Public uaa token must be PEM encoded")
	}

	return nil
}

func checkTokenFormat(token string) (string, error) {
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("Invalid token format")
	}

	tokenType, userToken := tokenParts[0], tokenParts[1]
	if !strings.EqualFold(tokenType, "bearer") {
		return "", errors.New("Invalid token type: " + tokenType)
	}

	return userToken, nil
}
