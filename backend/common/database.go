package common

import (
	"backend/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	host := viper.GetString("sql.host")
	port := viper.GetString("sql.port")
	database := viper.GetString("sql.database")
	username := viper.GetString("sql.username")
	password := viper.GetString("sql.password")
	loc := viper.GetString("sql.timezone")
	args := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s TimeZone=%s",
		host,
		port,
		username,
		database,
		password,
		loc)
	log.Println(args)
	db, err := gorm.Open(postgres.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	CreateTable(db)

	log.Printf("link success")
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

// 建表
func CreateTable(DB *gorm.DB) {
	DB.AutoMigrate(&model.Permission{})
	DB.AutoMigrate(&model.Tag{})
	DB.AutoMigrate(&model.Blog{})
	DB.AutoMigrate(&model.BlogTag{})
	DB.AutoMigrate(&model.New{})
	DB.AutoMigrate(&model.Problem{})
	DB.AutoMigrate(&model.Submission{})
	DB.AutoMigrate(&model.User{})
}
