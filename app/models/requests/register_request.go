package requests

type RegisterRequest struct {
	FirstName    string
	LastName     string
	Email           string
	Password        string
	ConfirmPassword string
}
