package models

import (
	"fiber_boilerplate/db"
	"strconv"

	"math/rand"

	"gorm.io/gorm"
)

type OTP struct {
	ID         int64  `gorm:"" db:"id, primarykey, autoincrement" json:"id"`
	UserID     int64  `gorm:"" db:"user_id" json:"user_id"`
	OtpCode    string `gorm:"" db:"otp_code" json:"otp_code" size:"6"`
	IsVerified int64  `gorm:"" db:"is_verified" json:"is_verified"`
	IsExpired  int16  `gorm:"" db:"is_expired" json:"is_expired"`
	table      string `gorm:"-"`
}

func (p OTP) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "otp"
}

type OtpModel struct{ DB *gorm.DB }

func (o *OtpModel) FindByUserID(userID int64) (*OTP, error) {
	var verify OTP
	if err := o.DB.Where("user_id = ? AND is_verified = 0 AND is_expired = 0", userID).First(&verify).Error; err != nil {
		return nil, err
	}
	return &verify, nil
}

func (o *OtpModel) GenerateRandomNumber() string {
	// Generate a random number between 100000 and 999999 (inclusive)
	num := rand.Intn(900000) + 100000
	return strconv.Itoa(num)
}

func (o OtpModel) Create(otp OTP) (u OTP, err error) {
	err = db.GetDB().Create(&otp).Error
	return otp, err
}

func (o OtpModel) FindOtp(userId int64, otp string) (u OTP, err error) {
	err = db.GetDB().Where("user_id = ?", userId).Where("otp_code = ? ", otp).Where("is_expired = 0").Where("is_verified = 0").First(&u).Error
	if err != nil {
		return OTP{}, err
	}
	return u, nil
}

func (o OtpModel) Update(otp OTP) (u OTP, err error) {
	err = db.GetDB().Updates(&otp).Error
	return otp, err
}

func (o OtpModel) Delete(id int64) (err error) {
	return db.GetDB().Where("id = ?", id).Delete(&OTP{}).Error
}
