package models

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostCreateRequest 创建文章请求
type PostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

// PostUpdateRequest 更新文章请求
type PostUpdateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}
