package handlers

import (
	"be2/internal/app"
	"be2/internal/domain"
	"be2/internal/http/v1/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const statusClientClosedRequest = 499

type Handler struct {
	AS     app.AddressService
	CS     app.ClientService
	logger *slog.Logger
}

func NewHandler(asi app.AddressService, csi app.ClientService, logger *slog.Logger) Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return Handler{
		AS:     asi,
		CS:     csi,
		logger: logger,
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
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	client, err := clientRow.ToDomainAddressClient()
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	clientID, err := h.CS.CreateClient(ctx, client)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	h.logger.InfoContext(ctx, "client created", "client_id", clientID)
	h.respondJSON(w, http.StatusCreated, dto.SuccessResponse{Status: "ok"})
}

func (h *Handler) HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var clientID dto.UUIDRequest
	if err := json.NewDecoder(r.Body).Decode(&clientID); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	id, err := clientID.ToDomain()
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	if deleted, err := h.CS.DeleteClient(ctx, id); err != nil {
		h.handleDomainError(w, err)
		return
	} else {
		h.logger.InfoContext(ctx, "client deleted", "client_id", id, "deleted_rows", deleted)
		h.respondJSON(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
	}
}

func (h *Handler) handleDomainError(w http.ResponseWriter, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "internal server error"
		details []dto.ErrorDetail
	)

	switch {
	case errors.Is(err, context.Canceled):
		status = statusClientClosedRequest
		message = "request canceled"
	case errors.Is(err, context.DeadlineExceeded):
		status = http.StatusGatewayTimeout
		message = "request deadline exceeded"
	}

	var validationErrs *domain.ValidationErrors
	switch {
	case errors.As(err, &validationErrs):
		status = http.StatusBadRequest
		message = "validation error"
		details = dto.FromValidationErrors(validationErrs)
	case errors.Is(err, domain.ErrClientAlreadyExists):
		status = http.StatusConflict
		message = "client already exists"
	case errors.Is(err, domain.ErrClientNotFound):
		status = http.StatusNotFound
		message = "client not found"
	case errors.Is(err, domain.ErrAddressNotFound):
		status = http.StatusNotFound
		message = "address not found"
	default:
		message = err.Error()
	}

	h.logger.Warn("request failed", "status", status, "message", message, "error", err)
	h.respondError(w, status, message, details)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string, details []dto.ErrorDetail) {
	h.respondJSON(w, status, dto.ErrorResponse{Error: message, Details: details})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}
