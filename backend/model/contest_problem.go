package model

import "gorm.io/gorm"

type ContestProblem struct {
	gorm.Model
	ContestId uint64 `json:"contest_id"`
	ProblemId uint64 `json:"problem_id"`
}
