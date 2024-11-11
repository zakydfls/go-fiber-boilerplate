package models

import (
	"fiber_boilerplate/db"
	"time"
)

type UserModel struct{}

type User struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Username      string    `db:"username" json:"username"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"-"`
	Phone         *string   `db:"phone" json:"phone"`
	Address       *string   `db:"address" json:"address"`
	Picture       *string   `db:"picture" json:"picture"`
	TwoFactorAuth bool      `db:"two_factor_auth" json:"two_factor_auth"`
	Role          string    `db:"role" json:"role"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	table         string
}

func (p *User) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "users"
}

func (m *UserModel) All() ([]User, int64, error) {
	var users *[]User
	err := db.GetDB().Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	count := int64(len(*users))
	return *users, count, nil
}

func (m *UserModel) Create(user *User) (*User, error) {
	err := db.GetDB().Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) FindByID(id int64) (*User, error) {
	var user *User
	if err := db.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) FindByEmail(email string) (*User, error) {
	var user *User
	if err := db.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) FindByUsername(username string) (*User, error) {
	var user *User
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) FindByPhone(phone string) (*User, error) {
	var user *User
	if err := db.GetDB().Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) Update(user *User) (*User, error) {
	err := db.GetDB().Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
