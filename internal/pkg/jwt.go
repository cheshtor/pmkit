package pkg

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	UserId int64
	jwt.RegisteredClaims
}

// JWT Token 签发秘钥
var signKey = "E~PAV)mQgcrh01ur&^6hf%_A(n67LPNaVVOf+yXL@VFmMb@SvrYfqwKarOxQhhKT"

// GenToken 生成 JWT Token
// userId 用户 ID
func GenToken(userId int64) (string, error) {
	claims := TokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        utils.UUID(),
			Issuer:    "pmkit",
			Subject:   fmt.Sprintf("User:%d", userId),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signKey))
	return tokenString, err
}

// ParseToken 解析 JWT Token 中携带的数据
func ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token invalid")
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("illegal token")
	}
	return claims, nil
}

// GetUserIdFromToken 从 JWT Token 中获取用户 ID
func GetUserIdFromToken(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return -1, err
	}
	return claims.UserId, nil
}

// GetRemainingMilliseconds 从 JWT Token 中获取 Token 剩余有效时间（单位：毫秒）
func GetRemainingMilliseconds(tokenString string) (int64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return -1, err
	}
	remainingTime := claims.ExpiresAt.Time.Sub(time.Now()).Milliseconds()
	return remainingTime, nil
}
