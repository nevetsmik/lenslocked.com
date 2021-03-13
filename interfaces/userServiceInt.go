package interfaces

import "lenslocked.com/models"

type UserServiceInt interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// You will receive either:
	// ErrNotFound, ErrInvalidPassword, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*models.User, error)
	// InitiateReset will complete all the model-related tasks
	// to start the password reset process for the user with
	// the provided email address. Once completed, it will
	// return the token, or an error if there was one.
	InitiateReset(email string) (string, error)
	// CompleteReset will complete all the model-related tasks
	// to complete the password reset process for the user that
	// the token matches, including updating that user's pw.
	// If the token has expired, or if it is invalid for any
	// other reason the ErrTokenInvalid error will be returned.
	CompleteReset(token, newPw string) (*models.User, error)
	UserDBInt
}
