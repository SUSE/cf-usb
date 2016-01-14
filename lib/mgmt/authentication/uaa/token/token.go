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

//go:generate counterfeiter -o fakes/fake_token.go . Token
type Token interface {
	DecodeToken(userToken string, desiredPermissions ...string) error
	CheckPublicToken() error
}

type NullToken struct{}

func (_ NullToken) DecodeToken(_ string, _ ...string) error {
	return nil
}

func (_ NullToken) CheckPublicToken() error {
	return nil
}

type accessToken struct {
	uaaPublicKey                string
	uaaSymmetricVerificationKey string
}

func NewAccessToken(uaaPublicKey string, uaaSymmetricVerificationKey string) accessToken {
	return accessToken{
		uaaPublicKey:                uaaPublicKey,
		uaaSymmetricVerificationKey: uaaSymmetricVerificationKey,
	}
}

func (accessToken accessToken) DecodeToken(userToken string, desiredPermissions ...string) error {
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

func (accessToken accessToken) CheckPublicToken() error {
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
