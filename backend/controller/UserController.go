package controller

import (
	"backend/common"
	"backend/dto"
	"backend/helper"
	"backend/model"
	"backend/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type updateUserRequest struct {
	UserName     string `json:"UserName"`
	UserNickname string `json:"UserNickname"`
	UserEmail    string `json:"UserEmail"`
	UserPassword string `json:"UserPassword"`
	NewPassword  string `json:"NewPassword"`
	IconUpload   string `json:"IconUpload"`
	UserIcon     string `json:"UserIcon"`
}

func Register(context *gin.Context) {
	// 获取参数  获取不到~
	//// 使用 map 获取请求的参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(context.Request.Body).Decode(&requestMap)
	//// 使用结构体来获取参数
	var requestUser = model.User{}
	//json.NewDecoder(context.Request.Body).Decode(&requestUser)
	// gin 框架提供的 bind 参数
	context.Bind(&requestUser)

	// 获取参数
	DB := common.GetDB()
	name := requestUser.UserName
	nickname := requestUser.UserNickname
	email := requestUser.UserEmail
	password := requestUser.UserPassword
	// 数据验证
	if len(name) == 0 || haveSpace(name) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "请输入正确的用户名")
		return
	}
	if len(email) == 0 || haveSpace(email) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "邮箱不能为空")
		return
	}
	/*
		邮箱验证 ： 待实现
	*/
	if len(password) < 6 || haveSpace(password) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能少于 6 位, 并且不能包含空格")
		return
	}
	if len(nickname) == 0 {
		nickname = helper.RandomString(6)
	}

	// 判断用户名是否存在
	if isUserNameExist(DB, name) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该用户已存在")
		return
	}
	// 创建用户 用户密码进行加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		UserName:     name,
		UserEmail:    email,
		UserNickname: nickname,
		UserPassword: string(hasedPassword),
		UserIcon:     "",
		PermissionId: 0,
	}
	DB.Create(&newUser)

	// 发放 token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}
	// 返回结果
	response.Success(context, gin.H{"token": token}, "登录成功")
}

func Login(context *gin.Context) {
	// 获取参数  获取不到~
	//// 使用 map 获取请求的参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(context.Request.Body).Decode(&requestMap)
	// 使用结构体来获取参数
	var requestUser = model.User{}
	//json.NewDecoder(context.Request.Body).Decode(&requestUser)
	// gin 框架提供的 bind 参数
	context.ShouldBind(&requestUser)
	// 获取参数
	DB := common.GetDB()
	name := requestUser.UserName
	password := requestUser.UserPassword
	log.Printf(name, password)
	// 数据验证
	if len(name) < 4 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户名必须大于 4 位")
		return
	}
	if len(password) < 6 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "密码不能少于 6 位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("user_name = ?", name).First(&user)
	if user.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		response.Response(context, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	// 发放 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}
	// 返回结果
	response.Success(context, gin.H{"token": token}, "登录成功")
}

func UpdateUserInfo(context *gin.Context) {
	var requestUser = updateUserRequest{}
	var UserInfo = model.User{}
	//json.NewDecoder(context.Request.Body).Decode(&requestUser)
	// gin 框架提供的 bind 参数
	context.ShouldBind(&requestUser)
	// 获取参数
	DB := common.GetDB()
	name := requestUser.UserName
	UserInfo = queryUserInfoByName(DB, name)
	nickname := requestUser.UserNickname
	email := requestUser.UserEmail
	password := requestUser.UserPassword
	newPassword := requestUser.NewPassword
	IconUpload := requestUser.IconUpload
	/* 更新昵称 */
	if len(nickname) != 0 && haveSpace(nickname) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户名长度为 4 - 12 个字符, 并且不能包含空格")
		return
	} else {
		UserInfo.UserNickname = nickname
	}
	/*
		更新密码 => 校验原始密码
	*/
	if len(newPassword) != 0 && haveSpace(newPassword) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "新密码长度为 6 - 12 个字符, 并且不能包含空格")
		return
	} else {
		if len(newPassword) != 0 {
			if len(newPassword) < 4 || len(newPassword) > 12 {
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "新密码长度为 6 - 12 个字符, 并且不能包含空格")
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(UserInfo.UserPassword), []byte(password)); err != nil {
				response.Response(context, http.StatusBadRequest, 400, nil, "原始密码错误")
				return
			}
			/* 密码加密 */
			hasedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				response.Response(context, http.StatusInternalServerError, 500, nil, "密码加密错误")
				return
			}
			UserInfo.UserPassword = string(hasedPassword)
		}
	}
	/* 邮箱验证 短信验证待实现 */
	if haveSpace(email) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "邮箱不能包含空格")
		return
	} else {
		UserInfo.UserEmail = email
	}
	if len(IconUpload) != 0 {
		/* 头像处理 */
		file, err := context.FormFile("UserIcon")
		if err != nil {
			// 处理获取文件出错的情况
			response.Fail(context, nil, "头像上传失败")
			return
		}
		MaxFileSize := 200 * 1024 // 200 kb
		fileSize := file.Size

		if fileSize > int64(MaxFileSize) {
			// 文件大小超过限制
			response.Fail(context, nil, "头像大小超过限制")
			return
		}

		/* 图像存储 -> 获取文件后缀名*/
		extension := filepath.Ext(file.Filename)

		newFileName := UserInfo.UserName + extension
		log.Printf(newFileName)
		/*保存在 static 目录下 */
		err = context.SaveUploadedFile(file, "./static/icon/"+newFileName)
		if err != nil {
			// 处理文件保存出错的情况
			response.Fail(context, nil, "文件存储失败")
			return
		}
		UserInfo.UserIcon = newFileName
		requestUser.UserIcon = newFileName
	} else {
		UserInfo.UserIcon = requestUser.UserIcon
	}
	// 返回结果
	DB.Model(&model.User{}).Where("id = ?", UserInfo.ID).Updates(UserInfo)
	response.Success(context, gin.H{"user": requestUser}, "用户信息修改成功")
}

// 获取用户信息
func Info(context *gin.Context) {
	user, _ := context.Get("user")
	context.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

// 查询用户是否存在
func isUserNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("user_name = ?", name).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// 查询字符串当中是否含有空格
func haveSpace(s string) bool {
	if strings.Contains(s, " ") {
		return true
	}
	return false
}

// 通过 user name 去数据库当中查询信息
func queryUserInfoByName(db *gorm.DB, name string) model.User {
	var user model.User
	db.Where("user_name = ?", name).First(&user)
	return user
}
