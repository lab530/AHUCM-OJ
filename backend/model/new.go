package model

import (
	"gorm.io/gorm"
)

type New struct {
	gorm.Model
	NewTitle   string `json:"new_title" gorm:"type:varchar(32); unique"`
	NewContext string `json:"new_context" gorm:"type:text;not null"`
	UserID     uint64 `json:"user_id"`
	Topping    bool   `json:"topping" gorm:"column:topping;type:boolean;default:0"` // 默认值为 0 不置顶
}
