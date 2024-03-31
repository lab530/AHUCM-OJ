package main

import (
	"backend/common"
	"backend/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func main() {
	InitConfig()
	common.InitDB()
	r := gin.Default()
	r = router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

// 获取 application 里面的配置
func InitConfig() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory: ", err)
	}

	// 拼接配置文件的相对路径
	configPath := filepath.Join(wd, "..", "config.toml")

	// 检查文件是否存在
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatal("Configuration file not found")
	}

	// 设置配置文件的路径
	viper.SetConfigFile(configPath)

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to read configuration file: ", err)
	}

	log.Println("Configuration file loaded successfully")
}
