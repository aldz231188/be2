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
	AS app.AddressService
	CS app.ClientService
}

func NewHandler(asi app.AddressService, csi app.ClientService) Handler {
	return Handler{
		AS: asi,
		CS: csi,
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
func (h *Handler) HandleCreateClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var clientRow dto.CreateClientRequest

	if err := json.NewDecoder(r.Body).Decode(&clientRow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	client, err := clientRow.ToDomainAddressClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.CS.CreateClient(ctx, client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok")

}

func (h *Handler) HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var clientId dto.UUIDRequest

	if err := json.NewDecoder(r.Body).Decode(&clientId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var mess string
	if id, err := clientId.ToDomain(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if _, err := h.CS.DeleteClient(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		mess = "-1"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ok" + mess)

}
