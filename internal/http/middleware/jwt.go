package middleware

import (
	"be2/internal/app"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userKey contextKey = "authUser"

type JWT struct {
	authService app.AuthService
}

func NewJWT(authService app.AuthService) *JWT {
	return &JWT{authService: authService}
}

func (m *JWT) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ValidateAccessToken(r.Context(), parts[1])
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func UserFromContext(ctx context.Context) string {
	if v := ctx.Value(userKey); v != nil {
		if user, ok := v.(string); ok {
			return user
		}
	}
	return ""
}
