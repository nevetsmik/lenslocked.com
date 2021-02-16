package interfaces

import "lenslocked.com/models"

type GalleryDBInt interface {
	Create(gallery *models.Gallery) error
}
