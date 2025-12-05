package http

import (
	"be2/internal/http/middleware"
	"be2/internal/http/v1/handlers"
	"context"
	"log"
	"net/http"

	_ "be2/swagger" // <- сгенерированные доки
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
)

func RegisterRoutes(lc fx.Lifecycle, h handlers.Handler, jwt *middleware.JWT) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", h.HandleLogin)
	mux.HandleFunc("POST /refresh", h.HandleRefresh)
	mux.HandleFunc("POST /logout", h.HandleLogout)
	mux.HandleFunc("POST /logout_all", h.HandleLogoutAll)
	mux.HandleFunc("POST /createclient", jwt.Protect(h.HandleCreateClient))
	mux.HandleFunc("POST /deleteclient", jwt.Protect(h.HandleDeleteClient))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler) // /swagger/index.html

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
