package http

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"
)

func RegisterRoutes(lc fx.Lifecycle, h Handler) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /addaddress", h.HandleAddAddress)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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

}
