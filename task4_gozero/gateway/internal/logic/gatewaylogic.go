package logic

import (
	"context"

	"gotask/task4_gozero/gateway/internal/svc"
	"gotask/task4_gozero/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayLogic {
	return &GatewayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewayLogic) Gateway(req *types.RegisterReq) (resp *types.LoginResp, err error) {
	return &types.LoginResp{Token: "ok", UserId: 0}, nil
}
