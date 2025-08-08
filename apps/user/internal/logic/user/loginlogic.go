package userlogic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"

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

var jwtSecret = "your-secret-key"

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	if in.Username == "admin" && in.Password == "123456" {
		uid := int64(1001)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"uid": uid,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenStr, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return nil, err
		}

		return &user.LoginResponse{
			Uid:   uid,
			Token: tokenStr,
		}, nil
	}
	return nil, errors.New("用户名或密码错误")
}
