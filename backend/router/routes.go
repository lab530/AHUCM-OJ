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

	router.POST("api/submit", controller.Submit)
	router.GET("api/submit/lang", controller.Getlang)
	router.GET("api/submit/list", controller.GetSubmitList)

	router.GET("api/admin/testcase", controller.GetTestCaseList)
	router.POST("api/admin/uploadcase", controller.UploadTestCase)
	router.DELETE("api/admin/deletecase", controller.DeleteTestCase)
	router.GET("api/admin/casedetail", controller.GetTestCaseDetail)
	router.PUT("api/admin/updatacase", controller.UpdateCase)
	return router
}
