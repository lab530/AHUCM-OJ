package dto

import "backend/model"

type UserDto struct {
	UserName     string `json:"user_name"`
	UserNickname string `json:"user_nickname"`
	UserPassword string `json:"user_password"`
	UserEmail    string `json:"user_email"`
	UserIcon     string `json:"user_icon"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		UserName:     user.UserName,
		UserNickname: user.UserNickname,
		UserEmail:    user.UserEmail,
		UserIcon:     user.UserIcon,
	}
}
