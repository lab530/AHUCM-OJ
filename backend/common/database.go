package common

import (
	"backend/migrations"
	"backend/model"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Tag{})
	DB.AutoMigrate(&model.Blog{})
	DB.AutoMigrate(&model.BlogTag{})
	DB.AutoMigrate(&model.New{})
	DB.AutoMigrate(&model.Problem{})
	DB.AutoMigrate(&model.Submission{})
	DB.AutoMigrate(&model.ProblemCategory{})
	DB.AutoMigrate(&model.Category{})
	DB.AutoMigrate(&model.Contest{})
	DB.AutoMigrate(&model.ContestProblem{})
	DB.AutoMigrate(&model.ContestUser{})
	DB.AutoMigrate(&model.ContestSubmit{})
	DB.AutoMigrate(&model.ContestRank{})
	if err := migrations.CreateTrigger(DB); err != nil {
		log.Println("触发器迁移失败")
		panic(err)
		return
	}
	CheckAndAddAdminUser(DB)
}

// 添加默认的 admin
func CheckAndAddAdminUser(db *gorm.DB) error {
	// 假设User模型有一个IsAdmin字段或者类似的字段来标识管理员
	var adminUser model.User
	// 查找IsAdmin为true的用户（或者根据其他条件查找）
	result := db.Where("user_name = ?", "admin").First(&adminUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果没有找到管理员用户，则添加一个新的
			Password := "adminn" // 替换为你的哈希密码
			/* 密码加密 */
			hasedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("密码加密错误")
				return err
			}
			admin := model.User{
				UserName:     "admin",
				UserNickname: "admin",
				UserEmail:    "admin",
				UserPassword: string(hasedPassword),
				PermissionId: 1,
				// ... 其他字段
			}
			if err := db.Create(&admin).Error; err != nil {
				return err // 返回错误
			}
			log.Println("Admin user created successfully")
		} else {
			return result.Error // 返回其他数据库错误
		}
	}
	return nil // 管理员用户已存在或成功添加
}
