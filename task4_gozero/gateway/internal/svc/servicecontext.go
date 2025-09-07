package svc

import (
	"gotask/task4_gozero/gateway/internal/config"
	"gotask/task4_gozero/rpc/blog/blogclient"
	"gotask/task4_gozero/rpc/user/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   userclient.User
	BlogRpc   blogclient.Blog
	JwtSecret string
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserRpc:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		BlogRpc:   blogclient.NewBlog(zrpc.MustNewClient(c.BlogRpc)),
		JwtSecret: c.JwtSecret,
	}
}
