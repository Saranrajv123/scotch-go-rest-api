package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);unique_index" json:"email"`
	FirstName    string `gorm:"size:100;not null" json:"firstname"`
	LastName     string `gorm:"size:100;not null" json:"lastname"`
	Password     string `gorm:"size:100;not null" json:"password"`
	ProfileImage string `gorm:"size:255" json:"profileimage"`
}

// hash the password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// check password hash and password from user input
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return errors.New("password incorrect")
	}

	return nil
}

func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedPassword, err := HashPassword(password)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// prepare strips user input of any white space

func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.ProfileImage = strings.TrimSpace(u.ProfileImage)
}

// validate the user input

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("Email is required")
		}

		if u.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	default:
		if u.FirstName == "" {
			return errors.New("FirstName is required")
		}

		if u.LastName == "" {
			return errors.New("LastName is required")
		}

		if u.Email == "" {
			return errors.New("Email is required")
		}

		if u.Password == "" {
			return errors.New("Password is required")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	}

}

// adds a user to the database
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, nil
	}

	return u, nil

}

// return user based on email
func (u *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}

	if err := db.Debug().Table("users").Where("email = ?", u.Email).First(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// return a list of all the user
func GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}

	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}

	return &users, nil
}
