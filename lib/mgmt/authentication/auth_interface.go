package authentication

//Authentication is the model to use for implementing authentication
type Authentication interface {
	IsAuthenticated(string) error
}
