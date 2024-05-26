package controller

import (
	"backend/common"
	"backend/dto"
	"backend/helper"
	"backend/model"
	"backend/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ContestAdd(context *gin.Context) {
	var NewContest = model.Contest{}
	StandTime := "2006-01-02T15:04:00"
	StartTime := context.PostForm("StartAt")
	EndTime := context.PostForm("EndAt")
	idString := context.PostForm("UserId")
	ProblemList := context.PostForm("ProblemList")
	Public := context.PostForm("Public")
	participants := context.PostForm("Participants")
	NewContest.Title = context.PostForm("Title")
	NewContest.Password = context.PostForm("Password")
	NewContest.Description = context.PostForm("Description")
	NewContest.Public = (Public == "1")

	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		// 处理转换错误
		response.Response(context, http.StatusInternalServerError, 500, nil, "用户 id 解析错误")
		return
	}
	NewContest.UserId = id
	log.Println(NewContest)
	st, err := time.Parse(StandTime, StartTime)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "初始时间解析错误")
		return
	}
	et, err := time.Parse(StandTime, EndTime)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "结束时间解析错误")
		return
	}

	NewContest.StartAt = st
	NewContest.EndAt = et
	DB := common.GetDB()

	// 返回结果
	if len(NewContest.Title) == 0 || helper.CheckAllEmptyString(NewContest.Title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "标题不能为空或者全是空格占位符")
		return
	}

	// 判断标题是否存在
	if isContestTitleExist(DB, NewContest.Title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该标题已存在")
		return
	}

	err = DB.Create(&NewContest).Error

	if err != nil {
		// 处理创建记录时的错误
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛添加错误")
		return
	}

	if !NewContest.Public {
		if len(participants) > 0 {
			participantNames := strings.Split(participants, ",")

			// 创建一个切片来存储有效的用户 ID
			userIDMap := make(map[int]bool)
			db := common.GetDB()
			// 遍历参与者列表
			for _, name := range participantNames {
				// 从数据库中查找用户 ID
				var user model.User
				result := db.Where("user_name = ?", strings.TrimSpace(name)).First(&user)
				if result.Error != nil {
					if result.Error == gorm.ErrRecordNotFound {
						// 如果用户不存在,跳过该用户
						continue
					} else {
						response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取私人竞赛名单错误")
						return
					}
				}
				userIDMap[int(user.ID)] = true
			}
			for userID := range userIDMap {
				contestUser := model.ContestUser{
					ContestId: uint64(NewContest.ID),
					UserId:    uint64(userID),
				}

				err := DB.Create(&contestUser).Error
				if err != nil {
					response.Response(context, http.StatusUnprocessableEntity, 422, nil, "添加参与者信息错误")
					return
				}
			}
		}
	}
	err = EditContestProblem(NewContest.ID, ProblemList)
	if err != nil {
		// 处理创建记录时的错误
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛题目添加错误")
		return
	}
	response.Success(context, nil, "竞赛添加成功")
	return
}

func GetContestList(context *gin.Context) {
	var ContestCount int64
	list := make([]*model.Contest, 0)
	tx := GetContest()
	err := tx.Omit("password").
		Order("id DESC").Find(&list).
		Count(&ContestCount).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}
	log.Println(list)
	response.Success(context, gin.H{
		"data": map[string]interface{}{
			"data":  list,
			"count": ContestCount,
		},
	}, "获取成功")
}
func ContestVerity(context *gin.Context) {
	DB := common.GetDB()
	cid := context.Query("cid")
	password := context.PostForm("ContestPassword")
	userid := context.PostForm("UserId")
	var ContestUser model.ContestUser
	contest := GetContestInfoByContestId(cid)
	if userid != "-1" {
		result := GetContestUser(userid, cid)
		if result == true {
			response.Success(context, nil, "您已通过验证")
			return
		}
		userid, err := strconv.ParseUint(userid, 10, 64)
		cid, err := strconv.ParseUint(cid, 10, 64)
		ContestUser.UserId = userid
		ContestUser.ContestId = cid
		if err != nil {
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "数据解析错误")
			return
		}
		if contest.ID != 0 {
			if contest.Public {
				if password != contest.Password {
					response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛密码验证错误")
					return
				} else {
					err := DB.Create(&ContestUser).Error
					if err != nil {
						response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛用户信息添加错误")
						return
					}
					response.Success(context, nil, "竞赛密码验证正确")
					return
				}
			} else {
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "本场比赛未公开或您暂未被邀请到本场比赛")
				return
			}
		}
	} else {
		if contest.Public == true {
			if password != contest.Password {
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛密码验证错误")
				return
			} else {
				response.Success(context, nil, "竞赛密码验证正确")
				return
			}
		} else {
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "本场比赛未公开或您暂未被邀请到本场比赛")
			return
		}
	}
	response.Success(context, nil, "竞赛密码验证正确")
}

func GetContestDetail(context *gin.Context) {
	db := common.GetDB()
	cid := context.Query("cid")
	var contest model.Contest
	if err := db.Where("id = ?", cid).Find(&contest).Error; err != nil {
		// 处理查询错误
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取竞赛信息错误")
		return
	}
	var ids []uint
	if err := db.Table("contest_problems").Where("contest_id = ? AND deleted_at IS NULL", contest.ID).Select("problem_id").Scan(&ids).Error; err != nil {
		// 处理查询错误
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取竞赛题目ID错误")
		return
	}
	adjustedIds := make([]string, len(ids))
	for i, id := range ids {
		adjustedIds[i] = strconv.Itoa(int(id) + 1000) // 将ID转换为字符串并加上1000
	}
	idString := strings.Join(adjustedIds, ",")
	if !contest.Public {
		var contestUsers []model.ContestUser
		if err := db.Where("contest_id = ?", contest.ID).Find(&contestUsers).Error; err != nil {
			// 处理查询错误
			response.Response(context, http.StatusUnprocessableEntity, 422, nil, "查询用户信息错误")
			return
		}
		userNames := make([]string, len(contestUsers))
		for i, user := range contestUsers {
			var userName string
			if err := db.Table("users").Where("id = ?", user.UserId).Pluck("user_name", &userName).Error; err != nil {
				// 处理查询错误
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户信息处理错误")
				return
			}
			userNames[i] = userName
		}
		namesString := strings.Join(userNames, ",")
		// 返回比赛信息和参与用户名称
		response.Success(context, gin.H{
			"data": gin.H{
				"data":     contest,
				"users":    namesString,
				"problems": idString,
			},
		}, "获取信息成功")
		return
	}
	response.Success(context, gin.H{"data": contest,
		"problems": idString,
	}, "获取信息成功")
}

func GetContestRankList(context *gin.Context) {
	DB := common.GetDB()
	cid, err := strconv.ParseUint(context.Query("cid"), 10, 64)
	if err != nil {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取竞赛错误")
		return
	}
	var rankList []struct {
		model.ContestRank
		UserName string `json:"name"`
	}
	DB.Model(&model.ContestRank{}).
		Where("contest_id = ?", cid).
		Joins("INNER JOIN users ON contest_ranks.user_id = users.id").
		Select("contest_ranks.*, users.user_name as user_name").
		Find(&rankList)
	response.Success(context, gin.H{"data": rankList}, "获取成功")
}

func GetContestSubmit(context *gin.Context) {
	cid := context.Query("cid")
	submitIds, err := GetSubmitIdByCid(cid)
	if err != nil {
		log.Println("Get Submit IDs Error:", err)
		response.Response(context, http.StatusInternalServerError, 422, nil, "获取 Submit ID 列表失败")
		return
	}
	DB := common.GetDB()
	// 根据 submitIds 查询 Submission 列表并包含用户信息
	var submissions []*model.Submission
	err = DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("user_password", "user_icon", "user_email", "user_nickname", "permission_id")
	}).Where("id IN ?", submitIds).Order("id DESC").Find(&submissions).Error
	if err != nil {
		log.Println("Get Submission List Error:", err)
		response.Response(context, http.StatusInternalServerError, 422, nil, "获取 Submission 列表失败")
		return
	}

	response.Success(context, gin.H{
		"data": submissions,
	}, "获取成功")
}

func UpdateContest(context *gin.Context) {
	var NewContest = model.Contest{}
	StandTime := "2006-01-02T15:04:00"
	StartTime := context.PostForm("StartAt")
	EndTime := context.PostForm("EndAt")
	ProblemList := context.PostForm("ProblemList")
	Public := context.PostForm("Public")
	participants := context.PostForm("Participants")
	NewContest.Title = context.PostForm("Title")
	NewContest.Password = context.PostForm("Password")
	NewContest.Description = context.PostForm("Description")
	NewContest.Public = (Public == "1")
	idStr := context.PostForm("ID")
	id64, err := strconv.ParseUint(idStr, 10, 64) // 第二个参数是基数（10 表示十进制），第三个参数是目标类型的位大小
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "获取竞赛 id 错误")
		return
	}
	NewContest.ID = uint(id64)
	log.Println(NewContest)
	st, err := time.Parse(StandTime, StartTime)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "初始时间解析错误")
		return
	}
	et, err := time.Parse(StandTime, EndTime)
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "结束时间解析错误")
		return
	}

	NewContest.StartAt = st
	NewContest.EndAt = et
	DB := common.GetDB()

	// 返回结果
	if len(NewContest.Title) == 0 || helper.CheckAllEmptyString(NewContest.Title) {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "标题不能为空或者全是空格占位符")
		return
	}
	// 判断标题是否存在
	var count int64
	result := DB.Model(&model.Contest{}).Where("title = ? AND id <> ?", NewContest.Title, NewContest.ID).Count(&count)
	if result.Error != nil {
		// 处理查询错误
		log.Printf("Error checking title existence: %v", result.Error)
		response.Response(context, http.StatusInternalServerError, 500, nil, "内部服务器错误")
		return
	}

	if count > 0 {
		// 标题已存在
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "该标题已存在")
		return
	}

	if err = DB.Model(&model.Contest{}).Where("id = ?", NewContest.ID).Updates(NewContest).Error; err != nil {
		// 处理更新错误
		log.Printf("Error updating contest: %v", err)
		// 返回错误响应或其他逻辑处理
		response.Response(context, http.StatusInternalServerError, 500, nil, "竞赛信息更新错误")
		return
	}

	if !NewContest.Public {
		if len(participants) > 0 {
			// 获取参赛用户名
			participantNames := strings.Split(participants, ",")
			// 创建一个map来存储用户名和ID的映射
			userIDsMap := make(map[string]uint)
			for _, name := range participantNames {
				name = strings.TrimSpace(name) // 去除用户名前后的空格
				var user model.User
				// 使用GORM的First方法查找第一个匹配的User记录，按Name字段进行查找
				result = DB.Where("user_name = ?", name).First(&user)
				// 不需要特别检查result.Error，因为GORM在找不到记录时不会返回错误
				// 如果找到了用户，将用户名和ID存储到map中
				if user.ID != 0 { // 或者你可以检查 user.Name 是否被设置了来确认是否找到了用户
					userIDsMap[user.UserName] = user.ID
				}
			}

			// 查询当前竞赛的所有参与者
			var currentContestUsers []model.ContestUser
			result = DB.Where("contest_id = ?", NewContest.ID).Find(&currentContestUsers)
			if result.Error != nil {
				response.Response(context, http.StatusUnprocessableEntity, 422, nil, "获取当前竞赛参与者列表错误")
				return
			}

			currentUserIDs := make(map[uint]bool) // 假设UserID是uint类型
			for _, user := range currentContestUsers {
				currentUserIDs[uint(user.UserId)] = true
			}

			// 遍历新名单，并检查用户是否在老名单中
			for _, name := range participantNames {
				name = strings.TrimSpace(name)
				var user model.User
				result := DB.Where("user_name = ?", name).First(&user)
				if result.Error != nil {
					// 处理查询错误（可能需要跳过或返回错误）
					continue
				}
				if user.ID != 0 { // 假设ID为0表示未找到用户
					// 检查用户是否已经在老名单中
					if _, exists := currentUserIDs[user.ID]; !exists {
						// 用户不在老名单中，需要添加到竞赛中
						newContestUser := model.ContestUser{
							ContestId: uint64(NewContest.ID), // 假设NewContest.ID是竞赛的ID
							UserId:    uint64(user.ID),
						}
						if err := DB.Create(&newContestUser).Error; err != nil {
							// 处理添加错误（可能需要记录日志或返回错误）
							response.Response(context, http.StatusUnprocessableEntity, 422, nil, "用户添加失败")
							continue
						}
					}
					// 无论是否添加，都从currentUserIDs中移除，以便后续检查删除
					delete(currentUserIDs, user.ID)
				}
			}

			// currentUserIDs中剩下的就是在老名单中但不在新名单中的用户ID
			for userID := range currentUserIDs {
				// 这里是删除逻辑，例如：
				var oldContestUser model.ContestUser
				result := DB.Where("contest_id = ? AND user_id = ?", NewContest.ID, userID).First(&oldContestUser)
				if result.Error == nil {
					if err := DB.Delete(&oldContestUser).Error; err != nil {
						// 处理删除错误
						response.Response(context, http.StatusUnprocessableEntity, 422, nil, "更新用户错误")
						continue
					}
				}
			}
		}
	}
	if len(ProblemList) > 0 {
		problemList := strings.Split(ProblemList, ",")
		realProblemIDs := make(map[uint]bool)
		for _, p := range problemList {
			id, err := strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Error parsing problem ID: %s\n", p)
				continue
			}
			realID := uint(id - 1000) // 减去1000得到真实ID
			realProblemIDs[realID] = true
		}

		var problemIDs []uint

		// 使用GORM查询contest_problems表中对应竞赛ID的所有ProblemID
		result = DB.Model(&model.ContestProblem{}).Where("contest_id = ?", NewContest.ID).Select("problem_id").Find(&problemIDs)
		if result.Error != nil {
			fmt.Printf("Error finding problem IDs: %v\n", result.Error)
			return
		}

		// 将 problemIDs 列表转换成 map 方便判断
		existingProblemIDs := make(map[uint]bool)
		for _, id := range problemIDs {
			existingProblemIDs[id] = true
		}

		// 找出新增的问题ID并添加到数据库
		for realID := range realProblemIDs {
			if !existingProblemIDs[realID] {
				// 立即添加新问题到数据库
				var problem model.Problem
				DB.Where("id = ?", realID).First(&problem)
				if problem.ID != 0 {
					newContestProblem := model.ContestProblem{
						ContestId: uint64(NewContest.ID),
						ProblemId: uint64(realID), // 使用realID
						// 设置其他字段（如果有的话）
					}
					if err = DB.Create(&newContestProblem).Error; err != nil {
						response.Response(context, http.StatusUnprocessableEntity, 422, nil, "更新题目错误")
						return
					}
				}
			}
		}
		for oldID := range existingProblemIDs {
			if !realProblemIDs[oldID] {
				// 尝试从数据库中删除该ContestProblem记录
				var contestProblem model.ContestProblem
				result := DB.Where("contest_id = ? AND problem_id = ?", NewContest.ID, oldID).First(&contestProblem)
				if result.Error == nil {
					if err := DB.Delete(&contestProblem).Error; err != nil {
						// 处理删除错误
						response.Response(context, http.StatusUnprocessableEntity, 422, nil, "更新题目错误")
						continue
					}
				}
			}
		}

	}
	response.Success(context, nil, "竞赛更新成功")
	return
}

func AuthContest(context *gin.Context) {
	uid := context.PostForm("UserId")
	cid := context.Query("cid")
	contest := GetContestInfoByContestId(cid)
	if contest.Public && len(contest.Password) == 0 {
		response.Success(context, nil, "验证成功")
		return
	}
	result := GetContestUser(uid, cid)
	if result {
		response.Success(context, nil, "竞赛密码验证正确")
		return
	} else {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "本场比赛未公开或您暂未被邀请到本场比赛, 若比赛公开请输入正确密码并重试")
		return
	}
}

// 获取竞赛提交的 submit id
func GetSubmitIdByCid(cid string) ([]uint64, error) {
	var submitIds []uint64
	db := common.GetDB()
	err := db.Model(&model.ContestSubmit{}).
		Where("contest_id = ?", cid).
		Pluck("submit_id", &submitIds).
		Error
	if err != nil {
		return nil, err
	}
	return submitIds, nil
}

// 获取指定 contest_id 下的所有 Problem 信息
func GetProblemsByContestId(context *gin.Context) {
	var problems []model.Problem
	db := common.GetDB()
	cid := context.Query("cid")
	// 1. 查询 ContestProblem 表,获取指定 contest_id 下的所有 problem_id
	// 2. 获取当前比赛的所有提交记录
	var contestProblems []model.ContestProblem
	result := db.Where("contest_id = ?", cid).Find(&contestProblems)
	if result.Error != nil {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "题目信息获取失败")
		return
	}
	// 2. 根据 problem_id 查询 Problem 表,获取全部问题信息
	var problemIds []uint64
	for _, cp := range contestProblems {
		problemIds = append(problemIds, cp.ProblemId)
	}
	db.Where("id IN (?)", problemIds).Find(&problems).Omit("simple_input", "simple_output", "data", "ProblemCategories")
	response.Success(context, gin.H{
		"data": problems,
	}, "竞赛密码验证正确")
	return
}

func GetContestInfo(context *gin.Context) {
	id := context.Query("cid")
	contest, err := getContestInfo(id)
	if err != nil {
		response.Response(context, http.StatusUnprocessableEntity, 422, nil, "竞赛信息获取错误")
		return
	}
	response.Success(context, gin.H{"data": contest}, "获取竞赛信息成功")
	return
}

// 通过关键字查找内容
func GetContest() *gorm.DB {
	DB := common.GetDB()
	tx := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("user_password", "user_icon", "user_email", "user_nickname", "permission_id")
	})
	return tx
}

func getContestInfo(id string) (dto.ContestDto, error) {
	db := common.GetDB()
	var contest model.Contest
	if err := db.Where("id = ?", id).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ContestDto{}, errors.New("Contest not found")
		}
		return dto.ContestDto{}, fmt.Errorf("Failed to fetch contest from the database: %w", err)
	}

	var user model.User
	if contest.UserId != 0 {
		if err := db.Where("id = ?", contest.UserId).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.ContestDto{}, errors.New("User not found")
			}
			return dto.ContestDto{}, fmt.Errorf("Failed to fetch user from the database: %w", err)
		}
	}

	return dto.ToContestDto(contest, user), nil
}

func isContestTitleExist(db *gorm.DB, title string) bool {
	var contest model.Contest
	db.Where("title = ?", title).First(&contest)
	if contest.ID != 0 {
		return true
	}
	return false
}

func EditContestProblem(id uint, ProblemList string) error {
	DB := common.GetDB()
	MaxProblemId := helper.GetLastProblemId(DB)
	if len(ProblemList) == 0 {
		return nil
	}
	if len(ProblemList) > 0 {
		problemIDs := strings.Split(ProblemList, ",")
		for _, problemIDStr := range problemIDs {
			problemID, err := strconv.Atoi(problemIDStr)
			if err != nil {
				return err
			}
			pid := uint(problemID - 1000)
			if pid <= MaxProblemId && pid > 0 {
				contestProblem := model.ContestProblem{
					ContestId: uint64(id),
					ProblemId: uint64(problemID) - 1000,
				}
				result := DB.Create(&contestProblem)
				if result.Error != nil {
					return result.Error
				}
			}

		}
	}
	return nil
}

func GetContestInfoByContestId(id string) model.Contest {
	DB := common.GetDB()
	var contest model.Contest
	DB.Where("id = ?", id).First(&contest)
	return contest
}

func GetContestUser(uid string, cid string) bool {
	DB := common.GetDB()
	var ContestUser model.ContestUser
	DB.Where("user_id = ? and contest_id = ?", uid, cid).First(&ContestUser)
	if ContestUser.ID != 0 {
		return true
	}
	return false
}
