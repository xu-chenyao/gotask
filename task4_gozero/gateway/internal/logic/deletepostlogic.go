package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gotask/task4_gozero/gateway/internal/svc"
)

type DeletePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePostLogic) DeletePost() error {
	// todo: add your logic here and delete this line

	return nil
}
