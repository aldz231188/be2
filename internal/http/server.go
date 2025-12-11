package http

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"
)

func NewServer(lc fx.Lifecycle, h *http.ServeMux) error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: h,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 15 * time.Second,
		// IdleTimeout:  60 * time.Second,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP Server on :8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP Server on :8080")
			return server.Shutdown(ctx)
		},
	})
	return nil

}
