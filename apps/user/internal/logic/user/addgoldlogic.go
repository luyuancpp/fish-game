package userlogic

import (
	"context"
	"strconv"

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
	uid := in.Uid
	amount := in.Amount

	uidInt64, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error("uid parse int error", err)
		return nil, err
	}

	// 查询用户
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, uidInt64)
	if err != nil {
		return nil, err
	}

	// 增加金币
	u.Gold += int64(amount)

	// 更新
	err = l.svcCtx.UserModel.Update(l.ctx, u)
	if err != nil {
		return nil, err
	}

	return &user.AddGoldResponse{
		Gold: int32(u.Gold),
	}, nil
}
