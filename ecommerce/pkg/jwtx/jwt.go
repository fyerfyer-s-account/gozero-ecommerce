package jwtx

import (
	"github.com/golang-jwt/jwt/v4"
)

// ClaimKey defines custom type for JWT claim keys
type ClaimKey string

// Define claim keys as constants
const (
    KeyUserId ClaimKey = "userId"
    KeyRole   ClaimKey = "role"
    KeyExp    ClaimKey = "exp"
    KeyIat    ClaimKey = "iat"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func GetToken(secretKey string, iat, seconds, userId int64, role string) (string, error) {
	claims := make(jwt.MapClaims)
    claims[string(KeyExp)] = iat + seconds
    claims[string(KeyIat)] = iat
    claims[string(KeyUserId)] = userId
    claims[string(KeyRole)] = role
    
    token := jwt.New(jwt.SigningMethodHS256)
    token.Claims = claims
    return token.SignedString([]byte(secretKey))
}
