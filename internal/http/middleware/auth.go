package middleware

import (
	"be2/internal/authz"
	"be2/internal/grpcutil"
	"go.uber.org/fx"
	"net/http"
	"strings"
)

type Auth struct {
	V *authz.Validator
	R *authz.RevokedSessions
}

func NewAuth(v *authz.Validator, r *authz.RevokedSessions) *Auth { return &Auth{V: v, R: r} }

func (m *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		claims, err := m.V.ParseAccess(strings.TrimSpace(h[len("bearer "):]))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		if m.R != nil && m.R.IsRevoked(claims.ID) {
			http.Error(w, "session revoked", http.StatusUnauthorized)
			return
		}
		ctx := grpcutil.WithUser(r.Context(), claims.UID, claims.Roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var Module = fx.Provide(NewAuth)
