package userlogic

import (
	"context"
	"errors"

	"fish-game/apps/user/fish-game/apps/user/user"
	"fish-game/apps/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if in.Username == "admin" && in.Password == "123456" {
		return &user.LoginResponse{
			Uid:   1001,
			Token: "mock-token-abc123",
		}, nil
	}
	return nil, errors.New("用户名或密码错误")
}
