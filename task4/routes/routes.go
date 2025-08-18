package routes

import (
	"task4/controllers"
	"task4/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 初始化控制器
	authController := &controllers.AuthController{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需token）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 文章相关路由
		posts := v1.Group("/posts")
		{
			// 公开路由（无需token）
			posts.GET("", postController.GetPosts)    // 获取文章列表
			posts.GET("/:id", postController.GetPost) // 获取单个文章

			// 需要认证的路由
			posts.POST("", middleware.AuthMiddleware(), postController.CreatePost)       // 创建文章
			posts.PUT("/:id", middleware.AuthMiddleware(), postController.UpdatePost)    // 更新文章
			posts.DELETE("/:id", middleware.AuthMiddleware(), postController.DeletePost) // 删除文章
		}

		// 评论相关路由
		comments := v1.Group("/comments")
		{
			// 公开路由
			comments.GET("/post/:post_id", commentController.GetCommentsByPost) // 获取文章评论

			// 需要认证的路由
			comments.POST("", middleware.AuthMiddleware(), commentController.CreateComment) // 创建评论
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	return r
}
