package roomlogic

import (
	"context"

	"fish-game/apps/room/fish-game/apps/room/room"
	"fish-game/apps/room/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveRoomLogic {
	return &LeaveRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveRoomLogic) LeaveRoom(in *room.LeaveRoomRequest) (*room.LeaveRoomReply, error) {
	// todo: add your logic here and delete this line

	return &room.LeaveRoomReply{}, nil
}
