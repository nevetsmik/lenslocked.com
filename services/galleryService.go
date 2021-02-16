package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/interfaces"
)

type galleryService struct {
	interfaces.GalleryDBInt
}

func NewGalleryService(db *gorm.DB) interfaces.GalleryServiceInt {
	return &galleryService{
		GalleryDBInt: &galleryValidator{
			GalleryDBInt: &galleryGorm{
				db: db,
			},
		},
	}
}
