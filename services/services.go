package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/interfaces"
	"lenslocked.com/models"
)

type Services struct {
	Gallery interfaces.GalleryServiceInt
	User    interfaces.UserServiceInt
	Image   interfaces.ImageServiceInt
	db      *gorm.DB
}

func NewServices(dialect, connectionInfo string) (*Services, error) {
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	// And next we need to construct services, but
	// we can't construct the UserService yet.
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&models.User{}, &models.Gallery{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&models.User{}, &models.Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
