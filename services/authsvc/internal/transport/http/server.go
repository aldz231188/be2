package httpserver

import (
	"context"
	"log"
	"net/http"

	"be2/services/authsvc/internal/config"
	"go.uber.org/fx"
)

func NewHTTPServer(cfg config.Config, jwks *JWKSHandler) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/.well-known/jwks.json", jwks)
	return &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}
}

func Run(lc fx.Lifecycle, srv *http.Server, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("[authsvc] jwks on %s", cfg.HTTPAddr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("[authsvc] jwks serve error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("[authsvc] jwks stopping")
			return srv.Shutdown(ctx)
		},
	})
}
