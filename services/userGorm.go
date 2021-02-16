package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/models"
)

// Implements UserDBInt

type UserGorm struct {
	db *gorm.DB
}

func (ug *UserGorm) ByID(id uint) (*models.User, error) {
	var user models.User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *UserGorm) ByEmail(email string) (*models.User, error) {
	var user models.User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *UserGorm) ByRemember(rememberHash string) (*models.User, error) {
	var user models.User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	return err
}

func (ug *UserGorm) Create(user *models.User) error {
	return ug.db.Create(user).Error
}

func (ug *UserGorm) Update(user *models.User) error {
	return ug.db.Save(user).Error
}

func (ug *UserGorm) Delete(id uint) error {
	user := models.User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

