package userlogic

import (
	"context"
	"fish-game/apps/user/user"

	"fish-game/apps/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProfileLogic) GetProfile(in *user.ProfileRequest) (*user.ProfileResponse, error) {
	// todo: add your logic here and delete this line

	return &user.ProfileResponse{}, nil
}
