package jwtx

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var DefaultSecret = "your-secret-key" // 可配置化

// GenerateToken 生成带 UID 的 JWT
func GenerateToken(uid int64, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(secret))
}

// VerifyToken 验证 JWT，返回 uid
func VerifyToken(tokenStr string, secret string) (int64, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if uidFloat, ok := claims["uid"].(float64); ok {
			return int64(uidFloat), true
		}
	}
	return 0, false
}
