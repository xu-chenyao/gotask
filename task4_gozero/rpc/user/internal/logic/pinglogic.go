package logic

import (
	"context"

	"gotask/task4_gozero/rpc/user/internal/svc"
	"gotask/task4_gozero/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	return &user.RegisterResponse{}, nil
}
