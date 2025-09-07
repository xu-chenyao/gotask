package logic

import (
	"context"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *blog.GetPostRequest) (*blog.GetPostResponse, error) {
	l.svcCtx.Mu.RLock()
	rec := l.svcCtx.Posts[in.Id]
	l.svcCtx.Mu.RUnlock()
	if rec == nil {
		return &blog.GetPostResponse{}, nil
	}
	return &blog.GetPostResponse{Post: &blog.Post{Id: rec.Id, Title: rec.Title, Content: rec.Content, UserId: rec.UserId, CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt}}, nil
}
