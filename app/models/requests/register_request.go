package requests

// RegisterRequest is a model for registering a new user.
type RegisterRequest struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	ConfirmPassword string
}
