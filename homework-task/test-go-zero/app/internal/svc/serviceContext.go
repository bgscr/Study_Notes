package svc

import (
	"blog-api/internal/config"
	"rpc/posts/postservice"
	"rpc/users/userservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRPC userservice.UserService
	PostRPC postservice.PostService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRPC: userservice.NewUserService(zrpc.MustNewClient(c.UserRPC)),
		PostRPC: postservice.NewPostService(zrpc.MustNewClient(c.PostRPC)),
	}
}
