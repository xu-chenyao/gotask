package logic

import (
	"context"
	"strings"

	"gotask/task4_gozero/gateway/internal/svc"
	"gotask/task4_gozero/gateway/internal/types"
	"gotask/task4_gozero/rpc/blog/blogclient"
	"gotask/task4_gozero/rpc/user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostReq) (resp *types.PostResp, err error) {
	// 从context取Authorization
	auth, _ := l.ctx.Value("Authorization").(string)
	var token string
	if strings.HasPrefix(auth, "Bearer ") {
		token = strings.TrimPrefix(auth, "Bearer ")
	}
	var userId int64
	if token != "" {
		vres, _ := l.svcCtx.UserRpc.Verify(l.ctx, &userclient.VerifyRequest{Token: token})
		if vres != nil && vres.Ok {
			userId = vres.UserId
		}
	}
	if userId == 0 { // 兜底
		userId = 1001
	}

	res, err := l.svcCtx.BlogRpc.CreatePost(l.ctx, &blogclient.CreatePostRequest{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userId,
	})
	if err != nil {
		return nil, err
	}
	p := res.Post
	if p == nil {
		return &types.PostResp{}, nil
	}
	return &types.PostResp{Id: p.Id, Title: p.Title, Content: p.Content, UserId: p.UserId, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}, nil
}
