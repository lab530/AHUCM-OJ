package model

type BlogTag struct {
	BlogTagId uint64 `json:"blog-tag-id" gorm:"type:int;primary_key"`
	BlogId    uint64 `json:"blog-id" gorm:"type:int;"`
	TagId     uint64 `json:"tag-id" gorm:"type:int;"`
}
