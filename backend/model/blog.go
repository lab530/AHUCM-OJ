package model

import "time"

type Blog struct {
	BlogId      uint64    `json:"blog-id" gorm:"type:int;primary_key"`
	BlogTitle   string    `json:"blog-title" gorm:"varchar(64);not null"`
	UserId      uint64    `json:"user-id" gorm:"type:int"`
	BlogContext string    `json:"blog-context" gorm:"text; not null"`
	UpdateTime  time.Time `json:"update-time" gorm:"timestamp"`
}
