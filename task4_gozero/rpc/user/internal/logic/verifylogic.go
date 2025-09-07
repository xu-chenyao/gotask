package logic

import (
	"context"
	"strconv"
	"strings"

	"gotask/task4_gozero/rpc/user/internal/svc"
	"gotask/task4_gozero/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLogic {
	return &VerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyLogic) Verify(in *user.VerifyRequest) (*user.VerifyResponse, error) {
	parts := strings.Split(in.Token, "|")
	if len(parts) != 2 {
		return &user.VerifyResponse{Ok: false}, nil
	}
	uid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return &user.VerifyResponse{Ok: false}, nil
	}
	// 简单校验：根据用户ID找到用户名（扫描内存）
	l.svcCtx.Mu.RLock()
	defer l.svcCtx.Mu.RUnlock()
	for _, u := range l.svcCtx.Users {
		if u.Id == uid {
			return &user.VerifyResponse{Ok: true, UserId: uid, Username: u.Username}, nil
		}
	}
	return &user.VerifyResponse{Ok: false}, nil
}
