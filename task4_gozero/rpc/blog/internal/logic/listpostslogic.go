package logic

import (
	"context"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *blog.ListPostsRequest) (*blog.ListPostsResponse, error) {
	l.svcCtx.Mu.RLock()
	defer l.svcCtx.Mu.RUnlock()
	var arr []*blog.Post
	for _, rec := range l.svcCtx.Posts {
		arr = append(arr, &blog.Post{Id: rec.Id, Title: rec.Title, Content: rec.Content, UserId: rec.UserId, CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt})
	}
	// 简化：不做严格分页
	return &blog.ListPostsResponse{Posts: arr}, nil
}
