package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/models"
)

// Implements UserDBInt

type userGorm struct {
	db *gorm.DB
}

func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

func (ug *userGorm) ByID(id uint) (*models.User, error) {
	var user models.User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*models.User, error) {
	var user models.User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByRemember(rememberHash string) (*models.User, error) {
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
		return ErrNotFound
	}
	return err
}

func (ug *userGorm) Create(user *models.User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *models.User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := models.User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Close() error {
	return ug.db.Close()
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}
