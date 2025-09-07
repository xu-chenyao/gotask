package logic

import (
	"context"
	"errors"

	"gotask/task4_gozero/rpc/user/internal/svc"
	"gotask/task4_gozero/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var userAutoId int64 = 1000

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errors.New("username/password required")
	}
	l.svcCtx.Mu.Lock()
	defer l.svcCtx.Mu.Unlock()
	if _, ok := l.svcCtx.Users[in.Username]; ok {
		return nil, errors.New("user exists")
	}
	userAutoId++
	rec := &svc.UserRecord{Id: userAutoId, Username: in.Username, Password: in.Password, Email: in.Email}
	l.svcCtx.Users[in.Username] = rec
	return &user.RegisterResponse{UserId: rec.Id}, nil
}
