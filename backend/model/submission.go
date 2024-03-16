package model

import (
	"backend/global"
	"github.com/jinzhu/gorm"
)

type Submission struct {
	gorm.Model                        // 包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` 字段
	UserId     uint64                 `json:"user_id"`
	Time       uint32                 `json:"time" gorm:"not null"` // 用时
	Memo       uint32                 `json:"memo" gorm:"not null"` // 所用内存
	Lang       string                 `json:"lang" gorm:"type:varchar(16);not null"`
	SourcePath string                 `json:"source_path" gorm:"type:varchar(128);not null"`
	Status     global.ExecutionResult `json:"status" gorm:"not null"`
}
