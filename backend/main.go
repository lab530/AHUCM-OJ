package main

import (
	"backend/common"
	"backend/helper"
	"backend/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	helper.InitConfig()
	common.InitDB()
	r := gin.Default()
	r = router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}
