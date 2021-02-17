package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/models"
)

// Implements GalleryDBInt interface

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) ByID(id uint) (*models.Gallery, error) {
	var gallery models.Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gg *galleryGorm) ByUserID (userID uint) ([]models.Gallery, error) {
	var galleries []models.Gallery
	db := gg.db.Where("user_id = ?", userID)
	if err := db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (gg *galleryGorm) Create(gallery *models.Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *models.Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := models.Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}
