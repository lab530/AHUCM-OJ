package model

type Tag struct {
	TagId    uint64 `json:"tag-id" gorm:"type:int;primary_key"`
	TageName string `json:"tage-name" gorm:"type:varchar(32);not null"`
}
