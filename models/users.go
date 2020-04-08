package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a recource cannot be found
	// in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: invalid password provided")
)

const userPwPepper = "lenslocked-secret-string"

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ByID will look up a user with the provided ID
// If the user is found, it will return a nil error
// If the user is not found, it will return ErrNotFound
// If there is any other error, it will return an error
// with more information
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id=?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address and
// returns that user.
// If the user is found, it will return a nil error
// If the user is not found, it will return ErrNotFound
// If there is any other error, it will return an error
// with more information
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

// Authenticate can be used to authenticate a user with the
// provided email address and password.
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, ErrNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// Update will update the provided users with all of the data
// in the provided User object.
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error

}

// Close closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops the User table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (us *UserService) AutoMigrate() error {
	return us.db.AutoMigrate(&User{}).Error
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}
