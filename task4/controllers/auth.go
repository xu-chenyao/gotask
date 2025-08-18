package controllers

import (
	"net/http"

	"task4/config"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthController 认证控制器
type AuthController struct{}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("用户注册参数验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.GetDB().Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "用户名或邮箱已存在",
		})
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: req.Password, // 密码将在BeforeCreate钩子中加密
		Email:    req.Email,
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		logrus.WithError(err).Error("创建用户失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败",
		})
		return
	}

	logrus.WithField("user_id", user.ID).Info("用户注册成功")
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("用户登录参数验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := config.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logrus.WithError(err).Error("生成token失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "登录失败",
		})
		return
	}

	logrus.WithField("user_id", user.ID).Info("用户登录成功")
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
