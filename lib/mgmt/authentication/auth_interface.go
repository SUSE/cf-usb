package authentication

type AuthenticationInterface interface {
	IsAuthenticated(string) error
}
