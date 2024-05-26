package model

import "gorm.io/gorm"

type ContestRank struct {
	gorm.Model
	ContestId    uint64 `json:"contest_id"`
	ProblemId    uint64 `json:"problem_id"`
	UserId       uint64 `json:"user_id"`
	Accepted     bool   `json:"accepted" gorm:"default:false"`
	SumSubmit    uint64 `json:"sum_submit"`
	PenaltyCount uint64 `json:"penalty_count"`
	Penalty      uint64 `json:"penalty"`
}
