package models

import (
	"fmt"
	"net/url"
	"path/filepath"
)

// Image is used to represent images stored in a Gallery.
// Image is NOT stored in the database, and instead references data stored on disk.
type Image struct {
	GalleryID uint
	Filename  string
}

// Path is used to build the absolute path used to reference this image via a web request.
func (i *Image) Path() string {
	// Create a properly encoded URL path - handles special characters in filenames
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local disk, relative to where our Go application is run from.
func (i *Image) RelativePath() string {
	// Convert the gallery ID to a string
	galleryID := fmt.Sprintf("%v", i.GalleryID)
	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}
