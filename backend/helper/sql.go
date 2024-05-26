package helper

import (
	"backend/model"
	"gorm.io/gorm"
)

func GetLastProblemId(db *gorm.DB) uint {
	problem := model.Problem{}
	result := db.Last(&problem)
	if result.Error != nil {
		panic("Failed to query last record")
	}
	return problem.ID
}
