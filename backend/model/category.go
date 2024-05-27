package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	CategoryName string `json:"category_name" gorm:"type:varchar(32);not null;unique"`
}
