package model

import "gorm.io/gorm"

type ContestSubmit struct {
	gorm.Model
	ContestId uint64 `json:"contest_id"`
	SubmitId  uint64 `json:"submit_id"`
}
