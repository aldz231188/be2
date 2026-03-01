package middleware

import (
	"be2/internal/clients/auth"
	"be2/internal/grpcutil"
	"go.uber.org/fx"
	"net/http"
	"strings"
)

type Auth struct{ C *auth.Service }

func NewAuth(c *auth.Service) *Auth { return &Auth{C: c} }

func (m *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(h[len("bearer "):])
		uid, err := m.C.ValidateAccess(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := grpcutil.WithUser(r.Context(), uid, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var Module = fx.Provide(NewAuth)
