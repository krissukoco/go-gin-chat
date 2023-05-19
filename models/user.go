package models

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AlphaNumerics = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrPasswordUnmatch = errors.New("password unmatch")
)

type User struct {
	Id        string `json:"id" gorm:"primaryKey"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	ImageUrl  string `json:"image_url"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func NewUserId() string {
	uid := "u_"
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		uid += string(AlphaNumerics[rand.Intn(len(AlphaNumerics))])
	}
	return uid
}

func (u *User) HashPassword() error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 6)
	if err != nil {
		return err
	}
	u.Password = string(hashedPwd)
	return nil
}

func (u *User) Save(db *gorm.DB) error {
	if u.Id == "" {
		u.Id = NewUserId()
	}
	tx := db.Save(&u)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (u *User) FindByUsername(db *gorm.DB, username string) error {
	tx := db.Where(&User{Username: username}).Take(&u)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return tx.Error
	}
	return nil
}

func (u *User) FindById(db *gorm.DB, id string) error {
	tx := db.Where(&User{Id: id}).Take(&u)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return tx.Error
	}
	return nil
}

func (u *User) ComparePassword(db *gorm.DB, rawPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPwd))
	if err != nil {
		return ErrPasswordUnmatch
	}
	return nil
}

func GetUsers(db *gorm.DB, offset int, limit int) ([]*User, error) {
	var users []*User
	tx := db.Offset(offset).
		Limit(limit).
		Order("created_at ASC").
		Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}
