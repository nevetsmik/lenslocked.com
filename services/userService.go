package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

// Top level service layer for users; implements userServiceInt

// Embeds a UserDB interface
type userService struct {
	interfaces.UserDBInt
}

var (
	ErrNotFound          = errors.New("models: resource not found")
	ErrIDInvalid         = errors.New("models: ID provided was invalid")
	ErrPasswordIncorrect = errors.New("models: incorrect password provided")
	ErrEmailRequired     = errors.New("models: email address is required")
	ErrEmailInvalid      = errors.New("models: email address is not valid")
	ErrEmailTaken        = errors.New("models: email address is already taken")
	ErrPasswordTooShort  = errors.New("models: password must be at least 8 characters long")
	ErrPasswordRequired  = errors.New("models: password is required")
	ErrRememberRequired  = errors.New("models: remember token is required")
	ErrRememberTooShort  = errors.New("models: remember token must be at least 32 bytes")
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
