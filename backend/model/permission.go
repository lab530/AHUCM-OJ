package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Privilege string `json:"privilege" gorm:"type:varchar(32); not null;"`
}
