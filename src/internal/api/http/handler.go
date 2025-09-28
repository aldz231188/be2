package http

import (
	// "be2/internal/app"
	"be2/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	DS domain.Service
}

func NewHandler(csi domain.Service) Handler {
	return Handler{
		DS: csi,
	}
}

func (h *Handler) HandleAddAddress(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var address Address

	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a := ToDomain(address)

	h.DS.AddAddress(ctx, a)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")

}
