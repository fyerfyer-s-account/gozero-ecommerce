package jwtx

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func GetToken(secretKey string, iat, seconds, userId int64, role string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["role"] = role
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
