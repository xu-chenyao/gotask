package controllers

import (
	"net/http"
	"strconv"

	"task4/config"
	"task4/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PostController 文章控制器
type PostController struct{}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req models.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("创建文章参数验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := config.GetDB().Create(&post).Error; err != nil {
		logrus.WithError(err).Error("创建文章失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建文章失败",
		})
		return
	}

	// 预加载用户信息
	config.GetDB().Preload("User").First(&post, post.ID)

	logrus.WithField("post_id", post.ID).Info("文章创建成功")
	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// GetPosts 获取文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// 查询文章列表
	if err := config.GetDB().
		Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		logrus.WithError(err).Error("获取文章列表失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文章列表失败",
		})
		return
	}

	// 获取总数
	var total int64
	config.GetDB().Model(&models.Post{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetPost 获取单个文章详情
func (pc *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文章ID",
		})
		return
	}

	var post models.Post
	if err := config.GetDB().
		Preload("User").
		Preload("Comments.User").
		First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文章ID",
		})
		return
	}

	var req models.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("更新文章参数验证失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数验证失败",
			"details": err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	// 查找文章
	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "只有文章作者才能更新文章",
		})
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content

	if err := config.GetDB().Save(&post).Error; err != nil {
		logrus.WithError(err).Error("更新文章失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新文章失败",
		})
		return
	}

	// 预加载用户信息
	config.GetDB().Preload("User").First(&post, post.ID)

	logrus.WithField("post_id", post.ID).Info("文章更新成功")
	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文章ID",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	// 查找文章
	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	// 检查是否为文章作者
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "只有文章作者才能删除文章",
		})
		return
	}

	// 删除文章（GORM会自动删除关联的评论）
	if err := config.GetDB().Delete(&post).Error; err != nil {
		logrus.WithError(err).Error("删除文章失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除文章失败",
		})
		return
	}

	logrus.WithField("post_id", post.ID).Info("文章删除成功")
	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
	})
}
