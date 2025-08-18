package controllers

import (
	"net/http"
	"strconv"

	"task4/config"
	"task4/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CommentController 评论控制器
type CommentController struct{}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("创建评论参数验证失败")
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

	// 检查文章是否存在
	var post models.Post
	if err := config.GetDB().First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  req.PostID,
	}

	if err := config.GetDB().Create(&comment).Error; err != nil {
		logrus.WithError(err).Error("创建评论失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建评论失败",
		})
		return
	}

	// 预加载用户信息
	config.GetDB().Preload("User").First(&comment, comment.ID)

	logrus.WithField("comment_id", comment.ID).Info("评论创建成功")
	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// GetCommentsByPost 获取文章的评论列表
func (cc *CommentController) GetCommentsByPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文章ID",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文章不存在",
		})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	// 查询评论列表
	var comments []models.Comment
	if err := config.GetDB().
		Where("post_id = ?", postID).
		Preload("User").
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error; err != nil {
		logrus.WithError(err).Error("获取评论列表失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取评论列表失败",
		})
		return
	}

	// 获取总数
	var total int64
	config.GetDB().Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}
