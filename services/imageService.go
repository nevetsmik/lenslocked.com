package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

type imageService struct{}

func NewImageService() interfaces.ImageServiceInt {
	return &imageService{}
}

// Create accepts the data to create an image - any type that satisfies the io.Reader interface, which is any type that
// has a Read method
// io.Copy is how we copy images and the src data is of type Reader
func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	path, err := is.mkImageDir(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]models.Image, error) {
	path := is.imageDir(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Setup the Image slice we are returning
	ret := make([]models.Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = models.Image{
			Filename:  filepath.Base(imgStr),
			GalleryID: galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) imageDir(galleryID uint) string {
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkImageDir(galleryID uint) (string, error) {
	// filepath.Join will return a path like: images/galleries/123
	galleryPath := is.imageDir(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

func (is *imageService) Delete(i *models.Image) error {
	return os.Remove(i.RelativePath())
}
