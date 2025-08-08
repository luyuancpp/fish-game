package roomlogic

import (
	"context"

	"fish-game/apps/room/fish-game/apps/room/room"
	"fish-game/apps/room/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPlayersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPlayersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPlayersLogic {
	return &GetPlayersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPlayersLogic) GetPlayers(in *room.GetPlayersRequest) (*room.GetPlayersReply, error) {
	// todo: add your logic here and delete this line

	return &room.GetPlayersReply{}, nil
}
