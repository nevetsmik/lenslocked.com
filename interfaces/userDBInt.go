package interfaces

import "lenslocked.com/models"

// UserDBInt is used to interact with the users database.
//
// For pretty much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error until we make "public"
// facing errors.

type UserDBInt interface {
	// Methods for querying for single users
	ByID(id uint) (*models.User, error)
	ByEmail(email string) (*models.User, error)
	ByRemember(token string) (*models.User, error)

	// Methods for altering users
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
}
