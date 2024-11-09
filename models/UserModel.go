package models

import "time"

type UserModel struct{}

type User struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Username      string    `db:"username" json:"username"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"-"`
	Phone         *string   `db:"phone" json:"phone"`
	Address       *string   `db:"address" json:"address"`
	Picture       string    `db:"picture" json:"picture"`
	TwoFactorAuth bool      `db:"two_factor_auth" json:"two_factor_auth"`
	Role          string    `db:"role" json:"role"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	table         string
}

func (p User) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "users"
}
