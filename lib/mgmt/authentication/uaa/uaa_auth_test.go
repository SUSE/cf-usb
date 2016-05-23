package uaa

import (
	"fmt"
	"testing"

	"github.com/pivotal-golang/lager/lagertest"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var signingKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1
JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R6
0gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB
AoGBAJ8dRTQFhIllbHx4GLbpTQsWXJ6w4hZvskJKCLM/o8R4n+0W45pQ1xEiYKdA
Z/DRcnjltylRImBD8XuLL8iYOQSZXNMb1h3g5/UGbUXLmCgQLOUUlnYt34QOQm+0
KvUqfMSFBbKMsYBAoQmNdTHBaz3dZa8ON9hh/f5TT8u0OWNRAkEA5opzsIXv+52J
duc1VGyX3SwlxiE2dStW8wZqGiuLH142n6MKnkLU4ctNLiclw6BZePXFZYIK+AkE
xQ+k16je5QJBAN0TIKMPWIbbHVr5rkdUqOyezlFFWYOwnMmw/BKa1d3zp54VP/P8
+5aQ2d4sMoKEOfdWH7UqMe3FszfYFvSu5KMCQFMYeFaaEEP7Jn8rGzfQ5HQd44ek
lQJqmq6CE2BXbY/i34FuvPcKU70HEEygY6Y9d8J3o6zQ0K9SYNu+pcXt4lkCQA3h
jJQQe5uEGJTExqed7jllQ0khFJzLMx0K6tj0NeeIzAaGCQz13oo2sCdeGRHO4aDh
HH6Qlq/6UOV5wP8+GAcCQFgRCcB+hrje8hfEEefHcFpyKH+5g1Eu1k0mLrxK2zd+
4SlotYRHgPCEubokb2S1zfZDWIXW3HmggnGgM949TlY=
-----END RSA PRIVATE KEY-----`

var uaaPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d
KVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX
qHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug
spULZVNRxq7veq/fzwIDAQAB
-----END PUBLIC KEY-----`

var wrongUaaPublicKey = `-----BEGIN PUBLIC KEY-----
no key
-----END PUBLIC KEY-----`

var symmetricKey = "x6ZyeOZG6rp80usI8M1J8AsisUE6v+QcU83ppNj3IVom0wWqj6hK9jUEec8noHDk"

var expiredToken = `bearer eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJmMDNhYzJhNS02OTBiLTQxNmMtOWU3Yy04YjVjYWE2ZDY0NWMiLCJzdWIiOiJjY191c2JfbWFuYWdlbWVudCIsImF1dGhvcml0aWVzIjpbInVzYi5tYW5hZ2VtZW50LmFkbWluIl0sInNjb3BlIjpbInVzYi5tYW5hZ2VtZW50LmFkbWluIl0sImNsaWVudF9pZCI6ImNjX3VzYl9tYW5hZ2VtZW50IiwiY2lkIjoiY2NfdXNiX21hbmFnZW1lbnQiLCJhenAiOiJjY191c2JfbWFuYWdlbWVudCIsImdyYW50X3R5cGUiOiJjbGllbnRfY3JlZGVudGlhbHMiLCJyZXZfc2lnIjoiNGRhMzk5NTMiLCJpYXQiOjE0NDQ5MTY0NDUsImV4cCI6MTQ0NDk1OTY0NSwiaXNzIjoiaHR0cHM6Ly91YWEuYm9zaC1saXRlLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjY191c2JfbWFuYWdlbWVudCIsInVzYi5tYW5hZ2VtZW50Il19.EJDHBommGHFAkDhErACG8EvxfW21cJEbqDeD8qOGD5Qzw7nIxISNxVgy3nn9henr52KQ03a68LkpWNPXYaLhdllr-x_dSH3pe-VNbc7mh12Fwfi-zLB0ILePf8jwD1xbABQ_p5QmUs-uwMKdwMHuvdlFR9opmcJpk9bMiTeTNGY`
var invalidToken = `bearer eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiTmMDNhYzJhNS02OTBiLTQxNmMtOWU3Yy04YjVjYWE2ZDY0NWMiLCJzdWIiOiJjY191c2JfbWFuYWdlbWVudCIsImF1dGhvcml0aWVzIjpbInVzYi5tYW5hZ2VtZW50LmFkbWluIl0sInNjb3BlIjpbInVzYi5tYW5hZ2VtZW50LmFkbWluIl0sImNsaWVudF9pZCI6ImNjX3VzYl9tYW5hZ2VtZW50IiwiY2lkIjoiY2NfdXNiX21hbmFnZW1lbnQiLCJhenAiOiJjY191c2JfbWFuYWdlbWVudCIsImdyYW50X3R5cGUiOiJjbGllbnRfY3JlZGVudGlhbHMiLCJyZXZfc2lnIjoiNGRhMzk5NTMiLCJpYXQiOjE0NDQ5MTY0NDUsImV4cCI6MTQ0NDk1OTY0NSwiaXNzIjoiaHR0cHM6Ly91YWEuYm9zaC1saXRlLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjY191c2JfbWFuYWdlbWVudCIsInVzYi5tYW5hZ2VtZW50Il19.EJDHBommGHFAkDhErACG8EvxfW21cJEbqDeD8qOGD5Qzw7nIxISNxVgy3nn9henr52KQ03a68LkpWNPXYaLhdllr-x_dSH3pe-VNbc7mh12Fwfi-zLB0ILePf8jwD1xbABQ_p5QmUs-uwMKdwMHuvdlFR9opmcJpk9bMiTeTNGY`
var logger = lagertest.NewTestLogger("mgmt-api")

func TestInitWrongUaaAuth(t *testing.T) {
	assert := assert.New(t)

	_, err := NewUaaAuth(wrongUaaPublicKey, "", "usb.management.admin", false, logger)
	assert.Error(err, "Public uaa token must be PEM encoded")
}

func TestDecodeExpiredToken(t *testing.T) {
	assert := assert.New(t)

	uaaauth, err := NewUaaAuth(uaaPublicKey, "", "usb.management.admin", false, logger)
	if err != nil {
		t.Errorf("Error initialising uaa auth: %v", err)
	}

	err = uaaauth.IsAuthenticated(expiredToken)
	assert.Error(err, "token is expired")
}

func TestDecodeInvalidToken(t *testing.T) {
	assert := assert.New(t)
	uaaauth, err := NewUaaAuth(uaaPublicKey, "", "usb.management.admin", false, logger)
	if err != nil {
		t.Errorf("Error initialising uaa auth: %v", err)
	}

	err = uaaauth.IsAuthenticated(invalidToken)
	assert.Error(err, "invalid character '$' looking for beginning of value")
}

func TestCodeDecodeToken(t *testing.T) {
	uaaauth, err := NewUaaAuth(uaaPublicKey, "", "usb.management.admin", false, logger)
	if err != nil {
		t.Errorf("Error initialising uaa auth: %v", err)
	}

	header := map[string]interface{}{
		"alg": "RS256",
	}

	alg := "RS256"
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header

	claims := map[string]interface{}{
		"exp":   3404281214,
		"scope": []string{"usb.management.admin"},
	}
	token.Claims = claims

	signedKey, err := token.SignedString([]byte(signingKey))
	if err != nil {
		t.Errorf("Error getting signed key: %v", err)
	}
	signedKey = "bearer " + signedKey

	err = uaaauth.IsAuthenticated(signedKey)
	if err != nil {
		t.Errorf("Error decoding token: %v", err)
	}
}

func TestCodeDecodeWrongScopeToken(t *testing.T) {
	testScope := "a.scope"
	uaaauth, err := NewUaaAuth(uaaPublicKey, "", testScope, false, logger)
	if err != nil {
		t.Errorf("Error initialising uaa auth: %v", err)
	}

	header := map[string]interface{}{
		"alg": "RS256",
	}

	alg := "RS256"
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header

	claims := map[string]interface{}{
		"exp":   3404281214,
		"scope": []string{"usb.management.admin"},
	}
	token.Claims = claims

	signedKey, err := token.SignedString([]byte(signingKey))
	if err != nil {
		t.Errorf("Error getting signed key: %v", err)
	}
	signedKey = "bearer " + signedKey

	err = uaaauth.IsAuthenticated(signedKey)
	assert.Error(t, err, fmt.Sprintf("Token does not have %v scope", testScope))
}

func TestSymmetricCodeDecodeToken(t *testing.T) {
	uaaauth, err := NewUaaAuth("", symmetricKey, "usb.management.admin", false, logger)
	if err != nil {
		t.Errorf("Error initialising uaa auth: %v", err)
	}

	header := map[string]interface{}{
		"alg": "HS256",
	}

	alg := "HS256"
	signingMethod := jwt.GetSigningMethod(alg)
	token := jwt.New(signingMethod)
	token.Header = header

	claims := map[string]interface{}{
		"exp":   3404281214,
		"scope": []string{"usb.management.admin"},
	}
	token.Claims = claims

	signedKey, err := token.SignedString([]byte(symmetricKey))
	if err != nil {
		t.Errorf("Error getting signed key: %v", err)
	}
	signedKey = "bearer " + signedKey

	err = uaaauth.IsAuthenticated(signedKey)
	if err != nil {
		t.Errorf("Error decoding token: %v", err)
	}
}
