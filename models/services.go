package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServiceConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServiceConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServiceConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServiceConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithGallery() ServiceConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServiceConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func WithOAuth() ServiceConfig {
	return func(s *Services) error {
		s.OAuth = NewOAuthService(s.db)
		return nil
	}
}

func NewServices(cfgs ...ServiceConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	Gallery GalleryService
	Image   ImageService
	User    UserService
	OAuth   OAuthService
	db      *gorm.DB
}

// Close closes the UserService database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops the User table and rebuilds it
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}, &OAuth{}, &pwReset{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &OAuth{}, &pwReset{}).Error
}
