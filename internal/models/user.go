package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string          `json:"email" gorm:"unique;not null"`
	Password       string          `json:"password,omitempty" gorm:"not null"`
	Phone          string          `json:"phone"`
	EmailVerified  bool            `json:"email_verified" gorm:"default:false"`
	PhoneVerified  bool            `json:"phone_verified" gorm:"default:false"`
	MembershipType MembershipLevel `json:"membership_type" gorm:"default:'BASIC'"`
	Name           string          `json:"name"`
}
