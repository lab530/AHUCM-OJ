package dto

import "backend/model"

type UserDto struct {
	ID           uint   `json:"id"`
	UserName     string `json:"user_name"`
	UserNickname string `json:"user_nickname"`
	UserPassword string `json:"user_password"`
	UserEmail    string `json:"user_email"`
	UserIcon     string `json:"user_icon"`
	PermissionId uint64 `json:"permission_id"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		ID:           user.ID,
		UserName:     user.UserName,
		UserNickname: user.UserNickname,
		UserEmail:    user.UserEmail,
		UserIcon:     user.UserIcon,
	}
}
