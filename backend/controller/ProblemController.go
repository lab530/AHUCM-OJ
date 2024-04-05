package controller

import (
	"backend/common"
	"backend/define"
	"backend/helper"
	"backend/model"
	"backend/response"
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
	err = tx.Count(&ProblemCount).Omit("description", "input", "output", "simple_input", "simple_output", "illustrate").Offset(page).Limit(size).Find(&list).Error

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

func ProblemAdd(context *gin.Context) {
	var DB = common.GetDB()
	var Problem = model.Problem{}
	context.ShouldBind(&Problem)
	time := Problem.TimeLimit
	memo := Problem.MemoLimit
	title := Problem.Title
	id := Problem.ID
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

	// 判断用户名是否存在
	if isTitleExist(DB, title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该标题已存在")
		return
	}

	newProblem := model.Problem{
		Title:        title,
		UserId:       int64(id),
		Description:  desc,
		Input:        input,
		Output:       oput,
		SimpleInput:  siput,
		SimpleOutput: soput,
		Illustrate:   Ill,
		Data:         "",
		TimeLimit:    time,
		MemoLimit:    memo,
	}
	err := DB.Create(&newProblem).Error
	if err != nil {
		// 处理创建记录时的错误
		panic(err)
	}
	// 返回结果
	response.Success(context, nil, "题目添加成功")
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

	response.Success(context, gin.H{
		"data": problem,
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
