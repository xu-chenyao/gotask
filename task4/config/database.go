package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// 数据库配置
const (
	DBUser     = "root"
	DBPassword = "123456"
	DBHost     = "localhost"
	DBPort     = "3306"
	DBName     = "blog_db"
)

// InitDB 初始化数据库连接
func InitDB() {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	var err error
	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	log.Println("数据库连接成功！")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
