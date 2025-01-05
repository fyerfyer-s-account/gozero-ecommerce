package middleware

import (
	"net/http"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
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
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := jwtx.ParseToken(parts[1], m.config.Auth.AccessSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Role != jwtx.RoleAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
