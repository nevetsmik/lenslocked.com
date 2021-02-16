package interfaces

// UserDBInt is implemented by userGorm
// UserServiceInt is implemented by userService
// GalleryDBInt => galleryGorm
// GalleryServiceInt => GalleryService

type GalleryServiceInt interface {
	GalleryDBInt
}
