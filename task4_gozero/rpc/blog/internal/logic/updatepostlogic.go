package logic

import (
	"context"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *blog.UpdatePostRequest) (*blog.UpdatePostResponse, error) {
	l.svcCtx.Mu.Lock()
	defer l.svcCtx.Mu.Unlock()
	rec := l.svcCtx.Posts[in.Id]
	if rec == nil || rec.UserId != in.UserId {
		return &blog.UpdatePostResponse{}, nil
	}
	if in.Title != "" {
		rec.Title = in.Title
	}
	if in.Content != "" {
		rec.Content = in.Content
	}
	rec.UpdatedAt = rec.UpdatedAt + 1
	return &blog.UpdatePostResponse{Post: &blog.Post{Id: rec.Id, Title: rec.Title, Content: rec.Content, UserId: rec.UserId, CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt}}, nil
}
