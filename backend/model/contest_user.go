package model

import "gorm.io/gorm"

type ContestUser struct {
	gorm.Model
	ContestId uint64 `json:"contest_id"`
	UserId    uint64 `json:"user_id"`
}
