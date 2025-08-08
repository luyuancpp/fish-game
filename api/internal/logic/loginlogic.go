package logic

import (
	"context"

	"fish-game/api/internal/svc"
	"fish-game/api/internal/types"
	"fish-game/apps/user/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.UserReply, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserReply{
		Uid:   rpcResp.Uid,
		Token: rpcResp.Token,
	}, nil
}
