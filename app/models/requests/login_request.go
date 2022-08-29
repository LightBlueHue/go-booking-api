// Package requests contains api request models.
package requests

// LoginRequest is a model for logging an existing user in.
type LoginRequest struct {
	Email    string
	Password string
}
