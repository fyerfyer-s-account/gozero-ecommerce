package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	config config.Config
}

func NewAuthMiddleware(config config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		// Get token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Split token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.config.Auth.AccessSecret), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userId, ok := claims[string(jwtx.KeyUserId)].(float64)
		if !ok {
			http.Error(w, "Invalid user id claim", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), jwtx.KeyUserId, int64(userId))
		r = r.WithContext(ctx)
		// Passthrough to next handler if need
		next(w, r)
	}
}
