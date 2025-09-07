package svc

import (
	"sync"

	"gotask/task4_gozero/rpc/blog/internal/config"
)

type ServiceContext struct {
	Config   config.Config
	Mu       sync.RWMutex
	Posts    map[int64]*PostRecord
	Comments map[int64][]*CommentRecord // by postId
}

type PostRecord struct {
	Id        int64
	Title     string
	Content   string
	UserId    int64
	CreatedAt int64
	UpdatedAt int64
}

type CommentRecord struct {
	Id        int64
	Content   string
	UserId    int64
	PostId    int64
	CreatedAt int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Posts:    make(map[int64]*PostRecord),
		Comments: make(map[int64][]*CommentRecord),
	}
}
