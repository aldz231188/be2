package server

import (
	"be2/internal/http/middleware"
	"be2/internal/http/v1/handlers"
	_ "be2/internal/swagger" // <- сгенерированные доки
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func RegisterRoutes(h handlers.Handler, jwt *middleware.JWT) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.HandleRegister)
	mux.HandleFunc("POST /login", h.HandleLogin)
	mux.HandleFunc("POST /refresh", h.HandleRefresh)
	mux.HandleFunc("POST /logout", h.HandleLogout)
	mux.HandleFunc("POST /logout_all", h.HandleLogoutAll)
	mux.HandleFunc("POST /createclient", jwt.Protect(h.HandleCreateClient))
	// mux.HandleFunc("POST /deleteclient", jwt.Protect(h.HandleDeleteClient))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler) // /swagger/index.html
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}
