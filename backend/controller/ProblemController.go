package controller

import (
	"backend/common"
	"backend/define"
	"backend/global"
	"backend/helper"
	"backend/model"
	"backend/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func GetProblemList(context *gin.Context) {
	size, _ := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem page Error:", err)
		return
	}
	page = (page - 1) * size
	var ProblemCount int64
	// page == 1 == > offset 0 数据库从 0 开始

	keyword := context.Query("keyword")
	category := context.Query("category")

	tx := GetProblem(keyword, category)

	list := make([]*model.Problem, 0)
	err = tx.Count(&ProblemCount).Omit("description", "input", "output", "simple_input", "simple_output", "illustrate").Order("id").Offset(page).Limit(size).Find(&list).Error

	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}
	response.Success(context, gin.H{
		"data": map[string]interface{}{
			"data":  list,
			"count": ProblemCount,
		},
	}, "获取成功")
}

func ProblemRevise(context *gin.Context) {
	pid := context.DefaultQuery("pid", "0")
	var DB = common.GetDB()
	var Problem = model.Problem{}
	context.ShouldBind(&Problem)
	time := Problem.TimeLimit
	memo := Problem.MemoLimit
	title := Problem.Title
	id := Problem.UserId
	desc := Problem.Description
	input := Problem.Input
	oput := Problem.Output
	siput := Problem.SimpleInput
	soput := Problem.SimpleOutput
	Ill := Problem.Illustrate
	if len(title) == 0 || helper.CheckAllEmptyString(title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "标题不能为空或者全是空格占位符")
		return
	}

	if time == 0 || memo == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "请输入合法的空间限制和内存限制")
		return
	}
	if len(desc) == 0 || helper.CheckAllEmptyString(desc) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目描述不能为空")
		return
	}

	// 判断标题是否存在
	if pid == "0" && isTitleExist(DB, title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该标题已存在")
		return
	}

	DataDirectoryName := helper.UniqueName()
	directoryPath := "./static/testcases/" + DataDirectoryName

	newProblem := model.Problem{
		Title:        title,
		UserId:       int64(id),
		Description:  desc,
		Input:        input,
		Output:       oput,
		SimpleInput:  siput,
		SimpleOutput: soput,
		Illustrate:   Ill,
		Data:         "backend/static/testcases/" + DataDirectoryName,
		TimeLimit:    time,
		MemoLimit:    memo,
	}

	if pid != "0" {
		data := GetPath(DB, pid)
		newProblem.Data = data
		DB.Model(&model.Problem{}).Where("id = ?", pid).Updates(newProblem)
		response.Success(context, nil, "题目编辑成功")
		return
	}

	err := DB.Create(&newProblem).Error
	if err != nil {
		// 处理创建记录时的错误
		log.Println(newProblem)
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目上传错误")
		panic(err)
	}
	err = helper.CreateDirectory(directoryPath)
	if err != nil {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目数据文件夹创建失败")
		return
	} else {
		fmt.Println("Directory created successfully!")
	}
	// 返回结果
	response.Success(context, nil, "题目添加成功")
}

func GetProblemPass(context *gin.Context) {
	db := common.GetDB()
	pid := context.Query("pid")
	cid := context.Query("cid")

	var statuses []struct {
		Status string `json:"value"`
		Count  int64  `json:"name"`
	}

	// 如果 cid 不为空,则从 contest_submit 表中查找
	if cid != "" {
		var submitIds []int64
		db.Model(&model.ContestSubmit{}).
			Where("contest_id = ?", cid).
			Pluck("submit_id", &submitIds)
		result := db.Model(&model.Submission{}).
			Where("problem_id = ? AND id IN (?)", pid, submitIds).
			Select("CAST(status AS CHAR) AS status, COUNT(*) AS count").
			Group("status").
			Scan(&statuses)

		if result.Error != nil {
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取题目数据失败")
			panic("failed to query submissions")
			return
		}
	} else {
		// 如果 cid 为空,则从 submission 表中查找
		result := db.Model(&model.Submission{}).
			Where("problem_id = ?", pid).
			Select("CAST(status AS CHAR) AS status, COUNT(*) AS count").
			Group("status").
			Scan(&statuses)

		if result.Error != nil {
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取题目数据失败")
			panic("failed to query submissions")
			return
		}
	}

	// 将数字状态码转换为对应的状态描述
	for i := range statuses {
		if status, ok := global.ExecutionResultToString[global.ExecutionResult(helper.StringToInt(statuses[i].Status))]; ok {
			statuses[i].Status = status
		} else {
			statuses[i].Status = "Unknown"
		}
	}

	response.Success(context, gin.H{"data": statuses}, "获取数据成功")
}

func GetProblemDetail(context *gin.Context) {
	db := common.GetDB()
	id := context.Query("pid")
	if id == "" {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取题目失败")
		return
	}
	var problem model.Problem
	log.Println(id)
	db.Where("id = ?", id).
		Preload("ProblemCategories").Preload("ProblemCategories.Category").First(&problem)
	if problem.ID == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该题目信息不存在")
		return
	}
	if problem.UserId == 0 {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户信息获取失败")
		return
	}
	Uinfo := QueryUserInfoById(db, int(problem.UserId))

	response.Success(context, gin.H{
		"data": problem,
		"Info": Uinfo,
	}, "获取题目成功")
}

// 查找标题
func isTitleExist(db *gorm.DB, title string) bool {
	var problem model.Problem
	db.Where("title = ?", title).First(&problem)
	if problem.ID != 0 {
		return true
	}
	return false
}

// 通过关键字查找内容
func GetProblem(keyword string, category string) *gorm.DB {
	DB := common.GetDB()
	tx := DB.Model(new(model.Problem)).Preload("ProblemCategories").Preload("ProblemCategories.Category").
		Where("title like ? or description like ?", "%"+keyword+"%", "%"+keyword+"%")
	log.Printf(category)
	if category != "" {
		tx.Joins("RIGHT JOIN problem_categories pc on pc.problem_id = problems.id").
			Where("pc.category_id = (SELECT c.id FROM categories c WHERE c.id = ?)", category)
	}
	return tx
}

func GetPath(db *gorm.DB, pid string) string {
	var problem model.Problem
	db.Where("id = ?", pid).First(&problem)
	if problem.ID != 0 {
		return problem.Data
	}
	return ""
}
