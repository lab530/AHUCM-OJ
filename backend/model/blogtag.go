package model

import "github.com/jinzhu/gorm"

type BlogTag struct {
	gorm.Model
	BlogId uint64 `json:"blog-id"`
	TagId  uint64 `json:"tag-id"`
}
