package userlogic

import (
	"context"

	"fish-game/apps/user/internal/svc"
	"fish-game/apps/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGoldLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGoldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoldLogic {
	return &AddGoldLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGoldLogic) AddGold(in *user.AddGoldRequest) (*user.AddGoldResponse, error) {
	// todo: add your logic here and delete this line

	return &user.AddGoldResponse{}, nil
}
