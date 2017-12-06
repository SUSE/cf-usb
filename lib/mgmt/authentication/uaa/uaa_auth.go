package uaa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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
func NewUaaAuth(uaaPublicKey, symmetricVerificationKey, scope, tokenURL string, devMode bool, logger lager.Logger) (authentication.Authentication, error) {
	var token accessToken.Token

	if devMode {
		token = accessToken.NullToken{}
	} else {
		if uaaPublicKey == "" && tokenURL != "" {
			// Fetch the public key from UAA
			uaaURL, err := url.Parse(tokenURL)
			if err != nil {
				logger.Error("initialize-uaa-parse-url", err)
				return nil, err
			}
			uaaURL.Path += "/token_key"
			resp, err := http.Get(uaaURL.String())
			if err != nil {
				logger.Error("initialize-uaa-fetch-token-key", err)
				return nil, err
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				err = fmt.Errorf("Got unexpected status %d (%s)", resp.StatusCode, resp.Status)
				logger.Error("initialize-uaa-fetch-token-key", err)
				return nil, err
			}
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error("initialize-uaa-read-token-key", err)
				return nil, err
			}
			var tokenKeyResponse struct {
				Value string `json:"value"`
			}
			if err = json.Unmarshal(responseBody, &tokenKeyResponse); err != nil {
				logger.Error("initialize-uaa-token-key-parse-json", err)
				return nil, err
			}
			uaaPublicKey = tokenKeyResponse.Value
		}
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
