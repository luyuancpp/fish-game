package userlogic

import (
	"context"
	"errors"
	"fish-game/apps/user/internal/svc"
	"fish-game/apps/user/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var jwtSecret = "secret123"

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

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if in.Username != "admin" || in.Password != "123456" {
		return nil, errors.New("用户名或密码错误")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": 1001,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Uid:   1001,
		Token: tokenStr,
	}, nil
}
