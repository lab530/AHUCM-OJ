package model

import "gorm.io/gorm"

type BlogTag struct {
	gorm.Model
	BlogId uint64 `json:"blog-id"`
	TagId  uint64 `json:"tag-id"`
}
