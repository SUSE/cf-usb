package token_test

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	authentication "github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa/token"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	var (
		accessToken authentication.Token

		signedKey       string
		UserPrivateKey  string
		UAAPublicKey    string
		UAASymmetricKey string

		token    *jwt.Token
		symToken *jwt.Token
		err      error
	)

	BeforeEach(func() {
		UserPrivateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1\nJZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R6\n0gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB\nAoGBAJ8dRTQFhIllbHx4GLbpTQsWXJ6w4hZvskJKCLM/o8R4n+0W45pQ1xEiYKdA\nZ/DRcnjltylRImBD8XuLL8iYOQSZXNMb1h3g5/UGbUXLmCgQLOUUlnYt34QOQm+0\nKvUqfMSFBbKMsYBAoQmNdTHBaz3dZa8ON9hh/f5TT8u0OWNRAkEA5opzsIXv+52J\nduc1VGyX3SwlxiE2dStW8wZqGiuLH142n6MKnkLU4ctNLiclw6BZePXFZYIK+AkE\nxQ+k16je5QJBAN0TIKMPWIbbHVr5rkdUqOyezlFFWYOwnMmw/BKa1d3zp54VP/P8\n+5aQ2d4sMoKEOfdWH7UqMe3FszfYFvSu5KMCQFMYeFaaEEP7Jn8rGzfQ5HQd44ek\nlQJqmq6CE2BXbY/i34FuvPcKU70HEEygY6Y9d8J3o6zQ0K9SYNu+pcXt4lkCQA3h\njJQQe5uEGJTExqed7jllQ0khFJzLMx0K6tj0NeeIzAaGCQz13oo2sCdeGRHO4aDh\nHH6Qlq/6UOV5wP8+GAcCQFgRCcB+hrje8hfEEefHcFpyKH+5g1Eu1k0mLrxK2zd+\n4SlotYRHgPCEubokb2S1zfZDWIXW3HmggnGgM949TlY=\n-----END RSA PRIVATE KEY-----"
		UAAPublicKey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d\nKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX\nqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug\nspULZVNRxq7veq/fzwIDAQAB\n-----END PUBLIC KEY-----"
		UAASymmetricKey = "symmetric-secret"

		token = jwt.New(jwt.GetSigningMethod("RS256"))
		symToken = jwt.New(jwt.SigningMethodHS256)

		accessToken = authentication.NewAccessToken(UAAPublicKey, UAASymmetricKey)
	})

	Describe(".DecodeToken", func() {
		Context("when the token is valid", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp":   3404281214,
					"scope": []string{"route.advertise"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error", func() {
				err := accessToken.DecodeToken("bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error if the token type string is capitalized", func() {
				err := accessToken.DecodeToken("Bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error if the token type string is uppercase", func() {
				err := accessToken.DecodeToken("BEARER "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the token is signed with HMAC and is valid", func() {
			Context("when the token is valid", func() {
				BeforeEach(func() {
					claims := map[string]interface{}{
						"exp":   3404281214,
						"scope": []string{"route.advertise"},
					}
					symToken.Claims = claims

					signedKey, err = symToken.SignedString([]byte(UAASymmetricKey))
					Expect(err).NotTo(HaveOccurred())
				})

				It("does not return an error", func() {
					err := accessToken.DecodeToken("bearer "+signedKey, "route.advertise")
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the token is signed with a different key", func() {
				BeforeEach(func() {
					claims := map[string]interface{}{
						"exp":   3404281214,
						"scope": []string{"route.advertise"},
					}
					symToken.Claims = claims

					signedKey, err = symToken.SignedString([]byte("another-secret"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error", func() {
					err := accessToken.DecodeToken("bearer "+signedKey, "route.advertise")
					Expect(err).To(HaveOccurred())
				})
			})

		})

		Context("when a token is not valid", func() {
			It("returns an error if the user token is not signed", func() {
				err = accessToken.DecodeToken("bearer not-a-signed-key", "not a permission")
				Expect(err).To(HaveOccurred())
			})

			It("returns an invalid token format when there is no token type", func() {
				err = accessToken.DecodeToken("has-no-token-type", "not a permission")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid token format"))
			})

			It("returns an invalid token type when type is not bearer", func() {
				err = accessToken.DecodeToken("basic some-auth", "not a permission")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid token type: basic"))
			})
		})

		Context("expired time", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp": time.Now().Unix() - 5,
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())

				signedKey = "bearer " + signedKey
			})

			It("returns an error if the token is expired", func() {
				err = accessToken.DecodeToken(signedKey, "route.advertise")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("token is expired"))
			})
		})

		Context("permissions", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp":   time.Now().Unix() + 50000000,
					"scope": []string{"route.foo"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())

				signedKey = "bearer " + signedKey
			})

			It("returns an error if the the user does not have requested permissions", func() {
				err = accessToken.DecodeToken(signedKey, "route.my-permissions", "some.other.scope")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Token does not have 'route.my-permissions', 'some.other.scope' scope"))
			})
		})
	})

	Describe(".CheckPublicToken", func() {
		BeforeEach(func() {
			accessToken = authentication.NewAccessToken("not a valid pem string", "")
		})

		It("returns an error if the public token is malformed", func() {
			err = accessToken.CheckPublicToken()
			Expect(err).To(HaveOccurred())
		})
	})
})
