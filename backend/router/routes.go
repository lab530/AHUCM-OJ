package router

import (
	"backend/controller"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(router *gin.Engine) *gin.Engine {
	router.Use(middleware.CORSMiddleware())
	router.POST("/api/auth/register", controller.Register)
	router.POST("/api/auth/login", controller.Login)
	router.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	router.PUT("/api/auth/edit", controller.UpdateUserInfo)

	router.POST("api/problem/add", controller.ProblemRevise)
	router.GET("api/problem/list", controller.GetProblemList)
	router.PUT("api/problem/edit", controller.ProblemRevise)
	router.GET("api/problem", controller.GetProblemDetail)
	router.GET("api/problempass", controller.GetProblemPass)

	router.POST("api/submit", controller.Submit)
	router.GET("api/submit/lang", controller.Getlang)
	router.GET("api/submit/list", controller.GetSubmitList)

	router.GET("api/admin/testCase", controller.GetTestCaseList)
	router.POST("api/admin/uploadCase", controller.UploadTestCase)
	router.DELETE("api/admin/deleteCase", controller.DeleteTestCase)
	router.GET("api/admin/caseDetail", controller.GetTestCaseDetail)
	router.PUT("api/admin/updateCase", controller.UpdateCase)

	//router.POST("api/admin/new/add", controller.NewRevise)

	router.POST("api/contest/add", controller.ContestAdd)
	router.GET("api/contestSet", controller.GetContestList)
	router.GET("api/contestInfo", controller.GetContestInfo)
	router.POST("api/contestverity", controller.ContestVerity)
	router.POST("api/authcontest", controller.AuthContest)
	router.GET("api/contestproblem", controller.GetProblemsByContestId)
	router.GET("api/contestsubmit", controller.GetContestSubmit)
	router.GET("api/contestrank", controller.GetContestRankList)
	router.GET("api/contestdetail", controller.GetContestDetail)
	router.PUT("api/updateContest", controller.UpdateContest)
	return router
}
