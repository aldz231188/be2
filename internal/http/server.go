package server

import (
	"context"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

func RegisterServer(lc fx.Lifecycle, h *http.ServeMux, shutdowner fx.Shutdowner) error {

	server := &http.Server{
		Addr:         ":8080",
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP Server on :8080")
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP server failed: %v", err)
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP Server on :8080")
			return server.Shutdown(ctx)
		},
	})
	return nil

}
