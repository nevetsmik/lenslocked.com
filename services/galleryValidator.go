package services

import (
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

type galleryValidator struct {
	interfaces.GalleryDBInt
}

type galleryValFn func(*models.Gallery) error

func (gv *galleryValidator) Create(gallery *models.Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDBInt.Create(gallery)
}


func runGalleryValFns(gallery *models.Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(g *models.Gallery) error {
	if g.UserID <= 0 {
		return models.ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *models.Gallery) error {
	if g.Title == "" {
		return models.ErrTitleRequired
	}
	return nil
}
