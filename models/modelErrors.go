package models

import "strings"

// Implements PublicErrorInt
// userValidatorService outputs ModelError type messages
type ModelError string

// Any type with an Error() method that returns a string implements the error interface
func (e ModelError) Error() string {
	return string(e)
}

// Nicely outputs the error messages
func (e ModelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

var (
	// users
	ErrNotFound          ModelError = "models: resource not found"
	ErrIDInvalid         ModelError = "models: ID provided was invalid"
	ErrPasswordIncorrect ModelError = "models: incorrect password provided"
	ErrEmailRequired     ModelError = "models: email address is required"
	ErrEmailInvalid      ModelError = "models: email address is not valid"
	ErrEmailTaken        ModelError = "models: email address is already taken"
	ErrPasswordTooShort  ModelError = "models: password must be at least 8 characters long"
	ErrPasswordRequired  ModelError = "models: password is required"
	ErrRememberRequired  ModelError = "models: remember token is required"
	ErrRememberTooShort  ModelError = "models: remember token must be at least 32 bytes"
	ErrUserIDRequired    ModelError = "models: user ID is required"
	// gallery
	ErrTitleRequired ModelError = "models: title is required"
)
