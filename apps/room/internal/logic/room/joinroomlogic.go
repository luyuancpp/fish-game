package roomlogic

import (
	"context"

	"fish-game/apps/room/fish-game/apps/room/room"
	"fish-game/apps/room/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinRoomLogic {
	return &JoinRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinRoomLogic) JoinRoom(in *room.JoinRoomRequest) (*room.JoinRoomReply, error) {
	// 简单分配房间，比如 uid 末尾是 0-4 -> room001
	//uid := in.Uid
	roomId := "room001" // 实际应 hash 分配
	return &room.JoinRoomReply{
		RoomId: roomId,
	}, nil
}
