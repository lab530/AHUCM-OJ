package model

import (
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	gorm.Model
	Title       string    `json:"title" gorm:"varchar(64);not null"`
	StartAt     time.Time `json:"start_at" gorm:"not null"`
	EndAt       time.Time `json:"end_at" gorm:"not null"`
	Description string    `json:"description" gorm:"text"`
	UserId      uint64    `json:"user_id" gorm:"type:int;"`
	Public      bool      `json:"public"`
	Password    string    `json:"password"`
	User        User      `json:"user" gorm:"foreignKey:UserId"`
}
