package server

import (
	"be2/internal/http/middleware"
	"be2/internal/http/v1/handlers"
	_ "be2/internal/swagger" // <- сгенерированные доки
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func RegisterRoutes(h handlers.Handler, authMW *middleware.Auth) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.Register)
	// mux.HandleFunc("POST /login", h.Login)
	// mux.HandleFunc("POST /refresh", h.Refresh)
	// mux.HandleFunc("POST /logout", h.Logout)
	// mux.HandleFunc("POST /logout_all", h.LogoutAll)

	mux.Handle("POST /createclient", authMW.Require(http.HandlerFunc(h.CreateClient)))
	// mux.HandleFunc("POST /createclient", authMW.Require(h.HandleCreateClient))
	// mux.HandleFunc("POST /deleteclient", jwt.Protect(h.HandleDeleteClient))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler) // /swagger/index.html
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}
