package model

type User struct {
	UserId       uint64 `json:"user_id" gorm:"type:int;primary_key"`
	UserName     string `json:"user_name" gorm:"type:varchar(255);not null;unique"`
	UserNickname string `json:"user_nickname" gorm:"type:varchar(255);not null"`
	UserPassword string `json:"user_password" gorm:"type:varchar(255);not null"`
	UserEmail    string `json:"user_email" gorm:"type:varchar(255);not null"`
	UserIcon     string `json:"user_icon" gorm:"type:varchar(255)"`
	PermissionId uint64 `json:"permission_id" gorm:"not null;default:0"`
}
