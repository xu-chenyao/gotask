package logic

import (
	"context"

	"gotask/task4_gozero/gateway/internal/svc"
	"gotask/task4_gozero/gateway/internal/types"
	"gotask/task4_gozero/rpc/blog/blogclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPostsLogic) ListPosts(req *types.ListPostsReq) (resp *types.ListPostsResp, err error) {
	res, err := l.svcCtx.BlogRpc.ListPosts(l.ctx, &blogclient.ListPostsRequest{Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	var arr []types.PostResp
	for _, p := range res.Posts {
		arr = append(arr, types.PostResp{Id: p.Id, Title: p.Title, Content: p.Content, UserId: p.UserId, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt})
	}
	return &types.ListPostsResp{Posts: arr}, nil
}
