package handlers

import (
	"be2/internal/app"
	"be2/internal/app/usecase"
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
	// AS     app.AddressService
	CS     usecase.ClientUsecase
	Auth   app.AuthService
	logger *slog.Logger
}

func NewHandler(csi usecase.ClientUsecase, auth app.AuthService, logger *slog.Logger) Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return Handler{
		// AS:     asi,
		CS:     csi,
		Auth:   auth,
		logger: logger,
	}
}

// HandleLogin godoc
// @Summary     Authenticate user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body dto.LoginRequest true "User credentials"
// @Success     200 {object} dto.TokenResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /login [post]
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var creds dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	tokens, err := h.Auth.Authenticate(ctx, creds.Username, creds.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidCredentials) {
			h.respondError(w, status, "invalid credentials", nil)
			return
		}
		h.respondError(w, status, err.Error(), nil)
		return
	}

	h.respondJSON(w, http.StatusOK, dto.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// HandleRegister godoc
// @Summary     Register a new user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body dto.RegisterRequest true "Registration payload"
// @Success     201 {object} dto.TokenResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /register [post]
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	tokens, err := h.Auth.Register(ctx, req.Username, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			status = http.StatusConflict
			h.respondError(w, status, "user already exists", nil)
			return
		case errors.Is(err, app.ErrInvalidCredentials):
			status = http.StatusBadRequest
			h.respondError(w, status, "invalid credentials", nil)
			return
		}
		h.respondError(w, status, err.Error(), nil)
		return
	}

	h.respondJSON(w, http.StatusCreated, dto.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// HandleRefresh godoc
// @Summary     Refresh access token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body dto.RefreshRequest true "Refresh token"
// @Success     200 {object} dto.TokenResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /refresh [post]
func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	tokens, err := h.Auth.Refresh(ctx, req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			h.respondError(w, status, "invalid or expired token", nil)
			return
		}
		h.respondError(w, status, err.Error(), nil)
		return
	}

	h.respondJSON(w, http.StatusOK, dto.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// HandleLogout godoc
// @Summary     Logout current session
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body dto.LogoutRequest true "Refresh token"
// @Success     200 {object} dto.SuccessResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /logout [post]
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.Auth.LogoutCurrent(ctx, req.RefreshToken); err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			h.respondError(w, status, "invalid or expired token", nil)
			return
		}
		h.respondError(w, status, err.Error(), nil)
		return
	}

	h.respondJSON(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
}

// HandleLogoutAll godoc
// @Summary     Logout from all sessions
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body dto.LogoutRequest true "Refresh token"
// @Success     200 {object} dto.SuccessResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /logout_all [post]
func (h *Handler) HandleLogoutAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.Auth.LogoutAll(ctx, req.RefreshToken); err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			h.respondError(w, status, "invalid or expired token", nil)
			return
		}
		h.respondError(w, status, err.Error(), nil)
		return
	}

	h.respondJSON(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
}

// HandleCreateClient godoc
// @Summary     Create client
// @Tags        clients
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateClientRequest true "Client payload"
// @Success     201 {object} dto.SuccessResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     409 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /createclient [post]
// @Security    BearerAuth
func (h *Handler) HandleCreateClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var clientRow dto.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&clientRow); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// client, err := clientRow.ToDomainAddressClient()
	// if err != nil {
	// 	h.handleDomainError(w, err)
	// 	return
	// }

	clientID, err := h.CS.Create(ctx, clientRow.UserID, clientRow.ClientName, clientRow.ClientSurname)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	h.logger.InfoContext(ctx, "client created", "client_id", clientID)
	h.respondJSON(w, http.StatusCreated, dto.SuccessResponse{Status: clientID})
}

// HandleDeleteClient godoc
// @Summary     Delete client
// @Tags        clients
// @Accept      json
// @Produce     json
// @Param       request body dto.UUIDRequest true "Client ID"
// @Success     200 {object} dto.SuccessResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     401 {object} dto.ErrorResponse
// @Failure     404 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Router      /deleteclient [post]
// @Security    BearerAuth
// func (h *Handler) HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
// 	defer cancel()

// 	var clientID dto.UUIDRequest
// 	if err := json.NewDecoder(r.Body).Decode(&clientID); err != nil {
// 		h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
// 		return
// 	}

// 	id, err := clientID.ToDomain()
// 	if err != nil {
// 		h.handleDomainError(w, err)
// 		return
// 	}

// 	if deleted, err := h.CS.DeleteClient(ctx, id); err != nil {
// 		h.handleDomainError(w, err)
// 		return
// 	} else {
// 		h.logger.InfoContext(ctx, "client deleted", "client_id", id, "deleted_rows", deleted)
// 		h.respondJSON(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
// 	}
// }

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
