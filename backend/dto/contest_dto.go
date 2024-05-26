package dto

import (
	"backend/model"
	"time"
)

type ContestDto struct {
	ID          uint      `json:"ID"`
	Title       string    `json:"title" gorm:"varchar(64);not null"`
	StartAt     time.Time `json:"start_at" gorm:"not null"`
	EndAt       time.Time `json:"end_at" gorm:"not null"`
	Description string    `json:"description" gorm:"text"`
	UserName    string    `json:"user_name"`
	Public      bool      `json:"public"`
	Password    string    `json:"password"`
}

func ToContestDto(contest model.Contest, user model.User) ContestDto {
	password := "False"
	if len(contest.Password) > 0 {
		password = "True"
	}
	return ContestDto{
		ID:          contest.ID,
		Title:       contest.Title,
		StartAt:     contest.StartAt,
		EndAt:       contest.EndAt,
		Description: contest.Description,
		UserName:    user.UserName,
		Public:      contest.Public,
		Password:    password,
	}
}
