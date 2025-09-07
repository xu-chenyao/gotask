package logic

import (
	"context"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCommentsLogic) ListComments(in *blog.ListCommentsRequest) (*blog.ListCommentsResponse, error) {
	l.svcCtx.Mu.RLock()
	defer l.svcCtx.Mu.RUnlock()
	var arr []*blog.Comment
	for _, rec := range l.svcCtx.Comments[in.PostId] {
		arr = append(arr, &blog.Comment{Id: rec.Id, Content: rec.Content, UserId: rec.UserId, PostId: rec.PostId, CreatedAt: rec.CreatedAt})
	}
	return &blog.ListCommentsResponse{Comments: arr}, nil
}
