package model

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	BlogTitle   string `json:"blog-title" gorm:"type:varchar(64);not null"`
	UserId      uint64 `json:"user-id"`
	BlogContext string `json:"blog-context" gorm:"type:text; not null"`
}
