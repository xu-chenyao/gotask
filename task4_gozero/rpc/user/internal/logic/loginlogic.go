package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gotask/task4_gozero/rpc/user/internal/svc"
	"gotask/task4_gozero/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	l.svcCtx.Mu.RLock()
	rec, ok := l.svcCtx.Users[in.Username]
	l.svcCtx.Mu.RUnlock()
	if !ok || rec.Password != in.Password {
		return nil, errors.New("invalid credentials")
	}
	// 简单token: userId|timestamp
	token := fmt.Sprintf("%d|%d", rec.Id, time.Now().Unix())
	return &user.LoginResponse{Token: token, UserId: rec.Id}, nil
}
