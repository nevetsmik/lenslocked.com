package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"

	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
	"lenslocked.com/rand"
)

type userValidator struct {
	// Methods from UserDBInt are forwarded to a userValidator type
	// Whatever instantiates the UserDBInt field is what uv.UserDBInt refers to
	// So to enable interface chaining from a uv to ug, the UserDBInt field will be instantiated to a ug
	interfaces.UserDBInt
	hmac       hash.HMAC
	pepper     string
	emailRegex *regexp.Regexp
}

type userValFn func(*models.User) error

func NewUserValidator(udb interfaces.UserDBInt, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDBInt: udb,
		hmac:      hmac,
		pepper:    pepper,
		emailRegex: regexp.MustCompile(
			`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) ByID(id uint) (*models.User, error) {
	// Validate the ID
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	// If it is valid, call the next method in the chain and
	// return its results.
	return uv.UserDBInt.ByID(id)
}

func (uv *userValidator) ByEmail(email string) (*models.User, error) {
	user := models.User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDBInt.ByEmail(user.Email)
}

func (uv *userValidator) ByRemember(token string) (*models.User, error) {
	user := models.User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDBInt.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *models.User) error {
	if err := runUserValFns(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return uv.UserDBInt.Create(user)
}

func (uv *userValidator) Update(user *models.User) error {
	if err := runUserValFns(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}

	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.UserDBInt.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user models.User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDBInt.Delete(id)
}

func runUserValFns(user *models.User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// Normalizes password to a hashed password
func (uv *userValidator) bcryptPassword(user *models.User) error {
	if user.Password == "" {
		// We DO NOT need to run this if the password
		// hasn't been changed.
		return nil
	}

	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// Normalizer - Hashes a remember token
func (uv *userValidator) hmacRemember(user *models.User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// Validate remember token is set.
// remember token needs to be set cus the remember token is what is set on the cookie.
// However, remember tokens are not persisted.
func (uv *userValidator) setRememberIfUnset(user *models.User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

// Validate a user id is greater than 0 on delete
func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *models.User) error {
		if user.ID <= n {
			return models.ErrIDInvalid
		}
		return nil
	})
}

// Normalize email addresses
func (uv *userValidator) normalizeEmail(user *models.User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// Validate an email address exists
func (uv *userValidator) requireEmail(user *models.User) error {
	if user.Email == "" {
		return models.ErrEmailRequired
	}
	return nil
}

// Validate email format
func (uv *userValidator) emailFormat(user *models.User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return models.ErrEmailInvalid
	}
	return nil
}

// Validate unique email address
func (uv *userValidator) emailIsAvail(user *models.User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == models.ErrNotFound {
		// Email address is available if we don't find
		// a user with that email address.
		return nil
	}
	// We can't continue our validation without a successful
	// query, so if we get any error other than ErrNotFound we
	// should return it.
	if err != nil {
		return err
	}

	// If we get here that means we found a user w/ this email
	// address, so we need to see if this is the same user we
	// are updating, or if we have a conflict.
	if user.ID != existing.ID {
		return models.ErrEmailTaken
	}
	return nil
}

// Validate length of password
func (uv *userValidator) passwordMinLength(user *models.User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return models.ErrPasswordTooShort
	}
	return nil
}

// Validate password is given by user
func (uv *userValidator) passwordRequired(user *models.User) error {
	if user.Password == "" {
		return models.ErrPasswordRequired
	}
	return nil
}

// Validate PasswordHash is always given a value
func (uv *userValidator) passwordHashRequired(user *models.User) error {
	if user.PasswordHash == "" {
		return models.ErrPasswordRequired
	}
	return nil
}

// Validate remember tokens are at least 32 bytes
func (uv *userValidator) rememberMinBytes(user *models.User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return models.ErrRememberTooShort
	}
	return nil
}

// Validate remember token exists
func (uv *userValidator) rememberHashRequired(user *models.User) error {
	if user.RememberHash == "" {
		return models.ErrRememberRequired
	}
	return nil
}
