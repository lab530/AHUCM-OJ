package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	TageName string `json:"tage_name" gorm:"type:varchar(32);not null"`
}
