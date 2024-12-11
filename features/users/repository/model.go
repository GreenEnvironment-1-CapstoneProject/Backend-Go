package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID           string `gorm:"primary_key;type:varchar(50);not null;column:id"`
	Username     string `gorm:"type:varchar(255);column:username;unique"`
	Password     string `gorm:"type:varchar(255);not null;column:password"`
	Name         string `gorm:"type:varchar(255);column:name"`
	Email        string `gorm:"type:varchar(255);not null;column:email;unique"`
	Address      string `gorm:"type:varchar(255);column:address"`
	Gender       string `gorm:"type:varchar(255);column:gender"`
	Phone        string `gorm:"type:varchar(255);column:phone"`
	Exp          int    `gorm:"type:int;not null;column:exp"`
	Coin         int    `gorm:"type:int;not nullcolumn:coin"`
	AvatarURL    string `gorm:"type:varchar(255);column:avatar_url"`
	IsMembership bool   `gorm:"type:boolean;column:is_membership;default:false"`
}

type VerifyOTP struct {
	*gorm.Model
	ID        string    `gorm:"primary_key;type:varchar(50);not null;column:id"`
	Email     string    `gorm:"type:varchar(255);not null;column:email"`
	OTP       string    `gorm:"type:varchar(255);not null;column:otp"`
	ExpiredAt time.Time `gorm:"not null;column:expired_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *VerifyOTP) TableName() string {
	return "verify_otp"
}