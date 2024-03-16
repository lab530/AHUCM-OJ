package model

import (
	"github.com/jinzhu/gorm"
)

type Problem struct {
	gorm.Model
	Title        string `json:"title" gorm:"type:varchar(32); unique"`
	UserId       int64  `json:"user_id"`
	Description  string `json:"description" gorm:"type:text;not null"`
	Input        string `json:"input" gorm:"type:text;not null"`
	Output       string `json:"output" gorm:"type:text;not null"`
	SimpleInput  string `json:"simple_input" gorm:"type:text"`
	SimpleOutput string `json:"simple_output" gorm:"type:text"`
	Illustrate   string `json:"illustrate" gorm:"type:text"`
	Data         string `json:"data" gorm:"type:varchar(128);not null"`
	TimeLimit    uint32 `json:"time_limit"`
	MemoLimit    uint32 `json:"memo_limit"`
}
