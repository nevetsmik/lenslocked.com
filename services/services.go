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

type ServicesConfig func(*Services) error

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&models.User{}, &models.Gallery{}, &models.PwReset{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&models.User{}, &models.Gallery{}, &models.PwReset{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
