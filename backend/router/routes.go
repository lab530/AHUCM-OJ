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
	return router
}
