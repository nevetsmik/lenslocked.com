package services

import (
	"golang.org/x/crypto/bcrypt"
	"strings"

	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

// Top level service layer for users; implements userServiceInt

// Embeds a UserDB interface
type userService struct {
	interfaces.UserDBInt
}

// Implements PublicErrorInt
// userValidatorService outputs modelError type messages
type modelError string

// Any type with an Error() method that returns a string implements the error interface
func (e modelError) Error() string {
	return string(e)
}

// Nicely outputs the error messages
func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

var (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: ID provided was invalid"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email address is not valid"
	ErrEmailTaken        modelError = "models: email address is already taken"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 characters long"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrRememberRequired  modelError = "models: remember token is required"
	ErrRememberTooShort  modelError = "models: remember token must be at least 32 bytes"
)

var userPwPepper = "secret-random-string"
var hmacSecretKey = "secret-hmac-key"

func NewUserService(connectionInfo string) (*userService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	// https://eli.thegreenplace.net/2020/embedding-in-go-part-3-interfaces-in-structs/
	// UserDBInt field for UserValidator is initialized to ug, a userGorm service (struct) that implements the UserDBInt interface.
	// userService embeds the UserDBInt interface and instantiates uv, a userValidator service (struct)
	// with UserDBInt forwarded methods from the UserDBInt: ug initialization
	//uv := &UserValidator{
	//	hmac:      hmac,
	//	UserDBInt: ug,
	//}
	uv := NewUserValidator(ug, hmac)
	return &userService{
		UserDBInt: uv,
	}, nil
}

func (us *userService) Authenticate(email, password string) (*models.User, error) {
	foundUser, err := us.UserDBInt.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}
