package model

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId  int64     `json:"problem_id" gorm:"type:int;"`          // 分类名称
	CategoryId int64     `json:"category_id" gorm:"type:int;"`         // 父级 ID
	Category   *Category `gorm:"foreignKey:category_id;references:id"` // 关联分类的基础信息表
}
