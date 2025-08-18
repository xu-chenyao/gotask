package main

import (
	"log"

	"task4/config"
	"task4/models"
	"task4/routes"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	// 设置日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 自动迁移数据库表
	if err := migrateDatabase(); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	logrus.Info("博客API服务器启动在端口 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// migrateDatabase 数据库迁移
func migrateDatabase() error {
	db := config.GetDB()

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)

	if err != nil {
		return err
	}

	logrus.Info("数据库表迁移完成")
	return nil
}
