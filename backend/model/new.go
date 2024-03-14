package model

import "time"

type New struct {
	NewId      uint64    `json:"new_id" gorm:"type:int; primary_key"`
	NewTitle   string    `json:"new_title" gorm:"type:varchar(64); unique"`
	NewContext string    `json:"new_context" gorm:"type:text;not null"`
	UserID     uint64    `json:"user_id" gorm:"type:int;"`
	UpdateTime time.Time `json:"update_time" gorm:"type:timestamp"`
	Topping    bool      `json:"topping" gorm:"column:topping;type:tinyint(1);default:0"`
}
