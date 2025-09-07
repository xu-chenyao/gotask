package logic

import (
	"context"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *blog.DeletePostRequest) (*blog.DeletePostResponse, error) {
	l.svcCtx.Mu.Lock()
	defer l.svcCtx.Mu.Unlock()
	rec := l.svcCtx.Posts[in.Id]
	if rec == nil || rec.UserId != in.UserId {
		return &blog.DeletePostResponse{Ok: false}, nil
	}
	delete(l.svcCtx.Posts, in.Id)
	return &blog.DeletePostResponse{Ok: true}, nil
}
