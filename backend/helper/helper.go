package helper

import (
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")
	result := make([]byte, n)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func CheckAllEmptyString(str string) bool {
	for _, char := range str {
		if char != ' ' {
			return false
		}
	}
	return true
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

// 创建唯一名字
func UniqueName() string {
	currentTime := time.Now()
	ms := currentTime.Round(time.Millisecond).Format(".000")
	NewFileName := currentTime.Format("20060102150405") + ms[1:]
	return NewFileName
}

// 创建唯一文件夹
func CreateDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

