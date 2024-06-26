package controller

import (
	"backend/common"
	"backend/define"
	"backend/helper"
	"backend/model"
	"backend/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
  	"bytes"
)

type SubmitRequest struct {
	UserId    uint64    `json:"user_id" `
	Lang      string    `json:"lang" `
	ProblemId uint64    `json:"problem_id" `
	Time      time.Time `json:"time" `
	Code      string    `json:"code" `
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
	languages := viper.GetStringSlice("common.languages")
	response.Response(context, http.StatusOK, 200, gin.H{
		"lang": languages,
	}, "获取成功")
}

func Submit(context *gin.Context) {
	cid, _ := strconv.Atoi(context.DefaultQuery("cid", define.DefaultCid))
	log.Println(cid)
	//pid, _ := strconv.Atoi(context.Query("pid"))
	var request SubmitRequest
	context.ShouldBind(&request)
	//if err := context.ShouldBind(&request); err != nil {
	//	response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目信息获取错误")
	//	return
	//}
	// 待补充
	var submit = model.Submission{}
	submit.UserId = request.UserId
	submit.SubmitTime = request.Time
	submit.ProblemId = request.ProblemId
	submit.Lang = request.Lang

	helper.InitConfig()
	languages := viper.GetStringSlice("common.languages")
	suffix := viper.GetStringSlice("common.suffix")
  os.MkdirAll("./static/code", 0755)

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
			submit.SourcePath = "./static/code/" + NewFileName
			err := SaveCodeToFile(request.Code, submit.SourcePath)
			if err != nil {
				// 处理保存文件出错的情况
				// 返回适当的错误响应
        log.Println("Code file save failed.");
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "代码保存失败")
				return
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
		// 检查 code 文件夹是否存在
		directoryPath := "./static/testcases/"
		err = helper.CreateDirectory(directoryPath)
		if err != nil {
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目数据文件夹创建失败")
			return
		} else {
			fmt.Println("Directory created successfully!")
		}
		submit.SourcePath = "./static/code/" + NewFileName
		err = context.SaveUploadedFile(file, submit.SourcePath)
		//extension := filepath.Ext(file.Filename)
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
	submit.TimeUsed = 0
	submit.MemoUsed = 0
	log.Println(submit)

	DB := common.GetDB()
	err = DB.Create(&submit).Error

	if err != nil {
		// 处理创建记录时的错误
		response.Response(context, http.StatusBadRequest, 400, nil, "提交失败")
		panic(err)
	}

  sourcePathAbs, _ := filepath.Abs(submit.SourcePath)
  testcasesPathAbs, _ := filepath.Abs("../" + submit.TestcasesPath)   // trick: make . as ..
  jsonBody := fmt.Sprintf(`{
"source_path": "%s",
"lang": "%s",
"problem_id": %d,
"mem_limit": %d,
"time_limit": %d,
"testcases_path": "%s",
"submission_id": %d
}`, sourcePathAbs, submit.Lang, submit.ProblemId, p.MemoLimit, p.TimeLimit / 1000, testcasesPathAbs, submit.ID)
  log.Println(jsonBody)

	coreHost := viper.GetStringSlice("core.host")[0]
	corePort := viper.GetStringSlice("core.port")[0]
  coreSubmitUrl := fmt.Sprintf("http://%s:%s/api/v1/submit", coreHost, corePort)
  resp, err := http.Post(coreSubmitUrl, "application/json", bytes.NewBuffer([]byte(jsonBody)))
  if err != nil {
      log.Println("core no response!");
  }
  defer resp.Body.Close()

	// 返回结果
	log.Println(cid)
	if cid != 0 {
		// 说明是竞赛提交， 需要在竞赛表中添加记录。
		var ContestSubmit model.ContestSubmit
		ContestSubmit.SubmitId = uint64(submit.ID)
		ContestSubmit.ContestId = uint64(cid)
		err = DB.Create(&ContestSubmit).Error
		if err != nil {
			// 处理创建记录时的错误
			response.Response(context, http.StatusBadRequest, 400, nil, "竞赛提交添加失败")
			panic(err)
		}

		var ContestRank model.ContestRank
		ContestRank.ContestId = uint64(cid)
		ContestRank.ProblemId = submit.ProblemId
		ContestRank.UserId = submit.UserId
		ContestRank.SumSubmit = 1
		ContestRank.Penalty = 0
		ContestRank.PenaltyCount = 0
		var existingRank model.ContestRank
		err = DB.Where("contest_id = ? AND problem_id = ? AND user_id = ?",
			ContestRank.ContestId, ContestRank.ProblemId, ContestRank.UserId).
			First(&existingRank).Error
		if existingRank.ID == 0 {
			// 不存在记录, 插入新记录
			err = DB.Create(&ContestRank).Error
			if err != nil {
				response.Response(context, http.StatusBadRequest, 400, nil, "竞赛统计添加失败")
				return
			}
		} else {
			// 存在记录，更新sum_submit字段
			existingRank.SumSubmit++
			err = DB.Save(&existingRank).Error
			if err != nil {
				response.Response(context, http.StatusBadRequest, 400, nil, "竞赛统计更新失败")
				return
			}
		}
	}
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
