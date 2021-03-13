package services

import (
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/jinzhu/gorm"

	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

// Top level service layer for users; implements userServiceInt

// Embeds a UserDB interface
type userService struct {
	interfaces.UserDBInt
	interfaces.PwResetDBInt
	pepper string
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) interfaces.UserServiceInt {
	ug := &UserGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	// https://eli.thegreenplace.net/2020/embedding-in-go-part-3-interfaces-in-structs/
	// UserDBInt field for UserValidator is initialized to ug, a UserGorm service (struct) that implements the UserDBInt interface.
	// userService embeds the UserDBInt interface and instantiates uv, a userValidator service (struct)
	// with UserDBInt forwarded methods from the UserDBInt: ug initialization
	//uv := &UserValidator{
	//	hmac:      hmac,
	//	UserDBInt: ug,
	//}
	uv := NewUserValidator(ug, hmac, pepper)
	return &userService{
		UserDBInt:    uv,
		PwResetDBInt: NewPwResetValidator(&pwResetGorm{db}, hmac),
		pepper:       pepper,
	}
}

func (us *userService) Authenticate(email, password string) (*models.User, error) {
	foundUser, err := us.UserDBInt.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+us.pepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, models.ErrPasswordIncorrect
	default:
		return nil, err
	}
}

func (us *userService) InitiateReset(email string) (string, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return "", err
	}
	pwr := models.PwReset{
		UserID: user.ID,
	}
	if err := us.PwResetDBInt.CreatePwResetToken(&pwr); err != nil {
		return "", err
	}
	return pwr.Token, nil
}

func (us *userService) CompleteReset(token, newPw string) (*models.User, error) {
	pwr, err := us.PwResetDBInt.ByToken(token)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, models.ErrTokenInvalid
		}
		return nil, err
	}
	// If the pwReset is over 12 hours old, it is invalid and we should return the ErrTokenInvalid error
	if time.Now().Sub(pwr.CreatedAt) > (12 * time.Hour) {
		return nil, models.ErrTokenInvalid
	}

	user, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}
	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}
	us.PwResetDBInt.DeletePwResetToken(pwr.ID)
	return user, nil
}
