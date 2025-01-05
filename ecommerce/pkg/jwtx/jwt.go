package jwtx

import (
	"fmt"
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

// CustomClaims holds JWT claims
type CustomClaims struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

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

// ParseToken parses and validates JWT token
func ParseToken(tokenString string, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
