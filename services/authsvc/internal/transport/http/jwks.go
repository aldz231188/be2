package httpserver

import (
	"encoding/json"
	"net/http"

	"be2/services/authsvc/internal/jwtkeys"
)

type JWKSHandler struct {
	key jwtkeys.JWK
}

type jwksResponse struct {
	Keys []jwtkeys.JWK `json:"keys"`
}

func NewJWKSHandler(rsaKey *jwtkeys.RSAKey) (*JWKSHandler, error) {
	jwk, err := rsaKey.JWK()
	if err != nil {
		return nil, err
	}
	return &JWKSHandler{key: jwk}, nil
}

func (h *JWKSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jwksResponse{Keys: []jwtkeys.JWK{h.key}})
}
