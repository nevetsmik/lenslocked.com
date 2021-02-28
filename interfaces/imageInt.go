package interfaces

import (
	"io"

	"lenslocked.com/models"
)

type ImageServiceInt interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]models.Image, error)
	Delete(i *models.Image) error
}
