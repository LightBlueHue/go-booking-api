package requests

type RegisterRequest struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}
