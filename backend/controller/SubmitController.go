package controller

import (
	"backend/common"
	"backend/define"
	"backend/helper"
	"backend/model"
	"backend/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type SubmitRequest struct {
	UserId    uint64    `json:"UserId" `
	Lang      string    `json:"Lang" `
	ProblemId uint64    `json:"ProblemId" `
	Time      time.Time `json:"Time" `
	Code      string    `json:"Code" `
}

func SaveCodeToFile(code string, filePath string) error {
	// 将代码字符串写入文件
	log.Println(filePath)
	err := os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Getlang(context *gin.Context) {
	helper.InitConfig()
	languages := viper.GetStringSlice("languages")
	response.Response(context, http.StatusOK, 200, gin.H{
		"lang": languages,
	}, "获取成功")
}

func Submit(context *gin.Context) {
	cid, _ := strconv.Atoi(context.DefaultQuery("cid", define.DefaultSize))
	//pid, _ := strconv.Atoi(context.Query("pid"))
	var request SubmitRequest
	context.ShouldBind(&request)
	//if err := context.ShouldBind(&request); err != nil {
	//	response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目信息获取错误")
	//	return
	//}
	// 待补充
	log.Println(request)
	if cid != 0 {
		// 说明是竞赛提交， 需要在竞赛表中添加记录。
	}
	var submit = model.Submission{}
	submit.UserId = request.UserId
	submit.SubmitTime = request.Time
	submit.ProblemId = request.ProblemId
	submit.Lang = request.Lang

	helper.InitConfig()
	languages := viper.GetStringSlice("languages")
	suffix := viper.GetStringSlice("suffix")
	if len(request.Lang) > 0 {
		// 如果提交的是 Code
		index := -1
		for i, language := range languages {
			if language == submit.Lang {
				index = i
				break
			}
		}
		if index != -1 {
			fileSuffix := suffix[index]
			// 获取毫秒
			ms := request.Time.Round(time.Millisecond).Format(".000")
			NewFileName := request.Time.Format("20060102150405") + ms[1:] + "." + fileSuffix
			submit.SourcePath = "backend/static/code/" + NewFileName
			err := SaveCodeToFile(request.Code, "./static/code/"+NewFileName)
			if err != nil {
				// 处理保存文件出错的情况
				// 返回适当的错误响应
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "代码保存失败")
			}

		}
	} else {
		// 如果提交的是文件
		file, err := context.FormFile("File")
		if err != nil {
			// 处理获取文件出错的情况
			response.Response(context, http.StatusBadRequest, 400, nil, "文件上传失败")
			return
		}
		// 获取文件后缀名
		extension := filepath.Ext(file.Filename)
		index := -1

		for i, suf := range suffix {
			if suf == extension[1:] {
				index = i
				break
			}
		}
		if index != -1 {
			submit.Lang = languages[index]
		} else {
			response.Response(context, http.StatusBadRequest, 400, nil, "文件获取失败, 暂不支持该文件后缀名的编程语言")
			return
		}
		ms := request.Time.Round(time.Millisecond).Format(".000")
		NewFileName := request.Time.Format("20060102150405") + ms[1:] + extension
		err = context.SaveUploadedFile(file, "./static/code/"+NewFileName)
		//extension := filepath.Ext(file.Filename)
		submit.SourcePath = "backend/static/code/" + NewFileName
		if err != nil {
			// 处理文件保存出错的情况
			//response.Fail(context, nil, "文件存储失败")
			response.Response(context, http.StatusBadRequest, 400, nil, "系统错误，文件存储失败")
			return
		}
	}
	p, err := GetProblemDetailsByID(submit.ProblemId)
	if err != nil {
		response.Response(context, http.StatusBadRequest, 400, nil, "获取问题详情错误")
		return
	}
	path, err := GetPathByProblemId(submit.ProblemId)
	if err != nil {
		response.Response(context, http.StatusBadRequest, 400, nil, "获取测试点路径错误")
		return
	}
	submit.TestcasesPath = path
	submit.TimeLimit = p.TimeLimit
	submit.MemoLimit = p.MemoLimit
	log.Println(submit)

	DB := common.GetDB()
	err = DB.Create(&submit).Error
	if err != nil {
		// 处理创建记录时的错误
		response.Response(context, http.StatusBadRequest, 400, nil, "提交失败")
		panic(err)
	}
	// 返回结果
	response.Success(context, nil, "提交成功")
}

func GetSubmitList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem page Error:", err)
		return
	}
	page = (page - 1) * size
	var SubmitCount int64
	// page == 1 == > offset 0 数据库从 0 开始
	tx := GetSubmit()

	list := make([]*model.Submission, 0)
	err = tx.Omit("source_path", "testcases_path").
		Order("id DESC").Find(&list).
		Count(&SubmitCount).Error

	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}
	response.Success(context, gin.H{
		"data": map[string]interface{}{
			"data":  list,
			"count": SubmitCount,
		},
	}, "获取成功")
}

func GetProblemDetailsByID(problemID uint64) (*model.Problem, error) {
	var problem model.Problem

	// 执行查询，并只选择时间和备注字段
	db := common.GetDB()
	err := db.Select("time_limit, memo_limit").Where("id = ?", problemID).First(&problem).Error
	if err != nil {
		// 处理查询错误
		return nil, err
	}

	return &problem, nil
}

func GetPathByProblemId(problemID uint64) (string, error) {
	var problem model.Problem
	db := common.GetDB()
	err := db.Where("id = ?", problemID).First(&problem).Error
	if err != nil {
		return "", err
	}
	return problem.Data, nil
}

func GetSubmit() *gorm.DB {
	DB := common.GetDB()
	tx := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("user_password", "user_icon", "user_email", "user_nickname", "permission_id")
	})
	return tx
}
