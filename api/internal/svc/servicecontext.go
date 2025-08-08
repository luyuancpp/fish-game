package svc

import (
	"fish-game/api/internal/config"
	"fish-game/apps/user/user"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
	}
}
