package logic

import (
	"context"
	"time"

	"gotask/task4_gozero/rpc/blog/blog"
	"gotask/task4_gozero/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

var postAutoId int64 = 2000

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *blog.CreatePostRequest) (*blog.CreatePostResponse, error) {
	l.svcCtx.Mu.Lock()
	defer l.svcCtx.Mu.Unlock()
	postAutoId++
	now := time.Now().Unix()
	rec := &svc.PostRecord{Id: postAutoId, Title: in.Title, Content: in.Content, UserId: in.UserId, CreatedAt: now, UpdatedAt: now}
	l.svcCtx.Posts[rec.Id] = rec
	return &blog.CreatePostResponse{Post: &blog.Post{Id: rec.Id, Title: rec.Title, Content: rec.Content, UserId: rec.UserId, CreatedAt: rec.CreatedAt, UpdatedAt: rec.UpdatedAt}}, nil
}
