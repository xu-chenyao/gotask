package svc

import (
	"sync"

	"gotask/task4_gozero/rpc/user/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	Users     map[string]*UserRecord // username -> record
	Mu        sync.RWMutex
}

type UserRecord struct {
	Id       int64
	Username string
	Password string // 明文示例（生产应使用哈希）
	Email    string
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Users:  make(map[string]*UserRecord),
	}
}
