package model

import (
	"backend/global"
	"time"
)

type Submission struct {
	SubmitId   uint64                 `json:"submit_id" gorm:"type=int;primaryKey"`
	ProblemId  uint64                 `json:"problem_id" gorm:"type:int;"`
	UserId     uint64                 `json:"user_id" gorm:"type:int;"`
	Time       uint32                 `json:"time" gorm:"type=int;not null"`
	Memo       uint32                 `json:"memo" gorm:"type=int;not null"`
	SubmitTime time.Time              `json:"submit_time" gorm:"type=timestamp;not null"`
	Lang       string                 `json:"lang" gorm:"type=varchar(255);not null"`
	SourcePath string                 `json:"source_path" gorm:"type=varchar(255);not null"`
	Status     global.ExecutionResult `json:"status" gorm:"not null"`
}
