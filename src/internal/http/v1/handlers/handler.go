package handlers

import (
	"be2/internal/app"
	// "be2/internal/domain"
	"be2/internal/http/v1/dto"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	DS app.Service
}

func NewHandler(csi app.Service) Handler {
	return Handler{
		DS: csi,
	}
}

// GetClient godoc
// @Summary     Get client by ID
// @Tags        clients
// @Produce     json
// @Param       id path string true "Client ID" format(uuid)
// @Success     200 {object} contracts.ClientResponse
// @Failure     404 {object} contracts.ErrorResponse
// @Router      /clients/{id} [get]
// @Security    BearerAuth
func (h *Handler) HandleCreateAddress(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var addressRow dto.CreateAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&addressRow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	address := addressRow.ToDomainAdress()

	h.DS.CreateAddress(ctx, address)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")

}
