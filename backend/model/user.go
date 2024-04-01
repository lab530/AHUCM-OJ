package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `json:"user_name" gorm:"type:varchar(32);not null;unique"`
	UserNickname string `json:"user_nickname" gorm:"type:varchar(32);not null"`
	UserEmail    string `json:"user_email" gorm:"type:varchar(64);not null"`
	UserIcon     string `json:"user_icon" gorm:"type:varchar(128)"`
	UserPassword string `json:"user_password" gorm:"type:varchar(64)"`
	PermissionId uint64 `json:"permission_id" gorm:"not null;default:0"`
}
