package model

import "gorm.io/gorm"

type BlogTag struct {
	gorm.Model
	BlogId uint64 `json:"blog_id"`
	TagId  uint64 `json:"tag_id"`
}
