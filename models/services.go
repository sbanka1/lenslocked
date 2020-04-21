package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// Close closes the UserService database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops the User table and rebuilds it
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
