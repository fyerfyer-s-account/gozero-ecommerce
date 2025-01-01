package middleware

import (
	"net/http"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
	"github.com/golang-jwt/jwt/v4"
)

type AdminAuthMiddleware struct {
	config config.Config
}

func NewAdminAuthMiddleware(config config.Config) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		config: config,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		// Get JWT token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Validate token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(m.config.AdminAuth.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check admin role
		claims := token.Claims.(jwt.MapClaims)
		role, ok := claims[m.config.AdminAuth.RoleKey].(string)
		if !ok || role != jwtx.RoleAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Passthrough to next handler if need
		next(w, r)
	}
}
