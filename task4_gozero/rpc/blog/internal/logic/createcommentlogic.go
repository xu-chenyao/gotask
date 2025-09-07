package logic

import (
	"context"
	"time"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

var commentAutoId int64 = 3000

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *blog.CreateCommentRequest) (*blog.CreateCommentResponse, error) {
	l.svcCtx.Mu.Lock()
	defer l.svcCtx.Mu.Unlock()
	commentAutoId++
	now := time.Now().Unix()
	rec := &svc.CommentRecord{Id: commentAutoId, Content: in.Content, UserId: in.UserId, PostId: in.PostId, CreatedAt: now}
	l.svcCtx.Comments[in.PostId] = append(l.svcCtx.Comments[in.PostId], rec)
	return &blog.CreateCommentResponse{Comment: &blog.Comment{Id: rec.Id, Content: rec.Content, UserId: rec.UserId, PostId: rec.PostId, CreatedAt: rec.CreatedAt}}, nil
}
