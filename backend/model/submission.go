package model

import (
	"backend/global"
	"gorm.io/gorm"
	"time"
)

type Submission struct {
	gorm.Model                           // 包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` 字段
	ProblemId     uint64                 `json:"problem_id" gorm:"type:int;"`
	UserId        uint64                 `json:"user_id" gorm:"type:int;"`
	TimeUsed      uint32                 `json:"time_used" gorm:"type=int;not null"`
	MemoUsed      uint32                 `json:"memo_used" gorm:"type=int;not null"`
	SubmitTime    time.Time              `json:"submit_time" gorm:"not null"`
	Lang          string                 `json:"lang" gorm:"type=varchar(255);not null"`
	SourcePath    string                 `json:"source_path" gorm:"type=varchar(255);not null"`
	TestcasesPath string                 `json:"testcases_path" gorm:"type=varchar(255);not null"`
	Status        global.ExecutionResult `json:"status" gorm:"not null"`
	User          User                   `json:"user" gorm:"foreignKey:UserId"`
}
