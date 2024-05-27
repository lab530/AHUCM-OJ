package model

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	BlogTitle   string `json:"blog_title" gorm:"type:varchar(64);not null"`
	UserId      uint64 `json:"user_id"`
	BlogContext string `json:"blog_context" gorm:"type:text; not null"`
}
