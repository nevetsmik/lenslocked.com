package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/models"
)

// Implements GalleryDBInt interface

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *models.Gallery) error {
	return gg.db.Create(gallery).Error
}
