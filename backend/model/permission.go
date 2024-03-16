package model

import "github.com/jinzhu/gorm"

type Permission struct {
	gorm.Model
	Privilege string `json:"privilege" gorm:"type:varchar(32); not null;"`
}
