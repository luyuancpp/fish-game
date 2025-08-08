package logic

import (
	"context"

	"fish-game/api/internal/svc"
	"fish-game/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile(req *types.ProfileReq) (resp *types.UserReply, err error) {
	// todo: add your logic here and delete this line

	return
}
