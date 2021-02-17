package interfaces

import "lenslocked.com/models"

type GalleryDBInt interface {
	ByID(id uint) (*models.Gallery, error)
	Create(gallery *models.Gallery) error
}
