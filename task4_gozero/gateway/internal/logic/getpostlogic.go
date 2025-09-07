package logic

import (
	"context"

	"gotask/task4_gozero/gateway/internal/svc"
	"gotask/task4_gozero/gateway/internal/types"
	"gotask/task4_gozero/rpc/blog/blogclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostLogic) GetPost(id int64) (resp *types.PostResp, err error) {
	res, err := l.svcCtx.BlogRpc.GetPost(l.ctx, &blogclient.GetPostRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if res.Post == nil {
		return &types.PostResp{}, nil
	}
	p := res.Post
	return &types.PostResp{Id: p.Id, Title: p.Title, Content: p.Content, UserId: p.UserId, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}, nil
}
