package services

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

// Top level service layer for users; implements userServiceInt

// Embeds a UserDB interface
type userService struct {
	interfaces.UserDBInt
	pepper string
}

func NewUserService(db *gorm.DB, pepper, hmacKey string) *userService {
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
		UserDBInt: uv,
		pepper:    pepper,
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
