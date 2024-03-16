package model

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model
	TageName string `json:"tage-name" gorm:"type:varchar(32);not null"`
}
