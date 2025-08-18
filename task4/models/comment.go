package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
}

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,min=1"`
	PostID  uint   `json:"post_id" binding:"required"`
}
