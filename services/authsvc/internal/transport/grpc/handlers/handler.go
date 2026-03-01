package handlers

import (
	authv1 "be2/contracts/gen/auth/v1"
	"be2/services/authsvc/internal/app"
	"be2/services/authsvc/internal/domain"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

// const statusClientClosedRequest = 499

type Handler struct {
	Auth   app.AuthService
	logger *slog.Logger
	authv1.UnimplementedAuthServiceServer
}

func NewHandler(auth app.AuthService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		Auth:   auth,
		logger: logger,
	}
}

// func (h *Handler) HandleLogin(ctx context.Context, r *authv1.LoginRequest) (*authv1.LoginResponse, error) {

// 	// var creds dto.LoginRequest
// 	// creds.Login = r.Login
// 	// creds.Password = r.Password

// 	tokens, err := h.Auth.Authenticate(ctx, r.Login, r.Password)
// 	if err != nil {
// 		// status := http.StatusUnauthorized
// 		if errors.Is(err, app.ErrInvalidCredentials) {
// 			// h.respondError(w, status, "invalid credentials", nil)
// 			return nil, errors.New("invalid credentials")
// 		}
// 		// h.respondError(w, status, err.Error(), nil)
// 		return nil, err
// 	}

// 	return &authv1.LoginResponse{
// 		Tokens: &authv1.TokenPair{
// 			AccessToken:      tokens.AccessToken,
// 			AccessExpiresAt:  tokens.AccessExpiresAt,
// 			RefreshToken:     tokens.RefreshToken,
// 			RefreshExpiresAt: tokens.RefreshExpiresAt,
// 			SessionId:        tokens.SessionId,
// 		},
// 	}, nil
// }

func (h *Handler) Register(ctx context.Context, r *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

	// var req dto.RegisterRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	tokens, err := h.Auth.Register(ctx, r.Login, r.Password)
	if err != nil {
		// status := http.StatusInternalServerError
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			// status = http.StatusConflict
			// h.respondError(w, status, "user already exists", nil)
			return nil, errors.New("user already exists")
		case errors.Is(err, app.ErrInvalidCredentials):
			// status = http.StatusBadRequest
			// h.respondError(w, status, "invalid credentials", nil)
			return nil, errors.New("invalid credentials")
		}
		// h.respondError(w, status, err.Error(), nil)
		return nil, err
	}

	// h.respond(w, http.StatusCreated, dto.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
	return &authv1.RegisterResponse{
		Tokens: &authv1.TokenPair{
			AccessToken:      tokens.AccessToken,
			AccessExpiresAt:  tokens.AccessExpiresAt,
			RefreshToken:     tokens.RefreshToken,
			RefreshExpiresAt: tokens.RefreshExpiresAt,
			SessionId:        tokens.SessionId,
		},
	}, nil
}

func (h *Handler) HandleRefresh(ctx context.Context, r *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

	// var req dto.RefreshRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	tokens, err := h.Auth.Refresh(ctx, r.RefreshToken)
	if err != nil {
		// status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			// h.respondError(w, status, "invalid or expired token", nil)
			return nil, errors.New("invalid or expired token")
		}
		// h.respondError(w, status, err.Error(), nil)
		return nil, err
	}

	// h.respond(w, http.StatusOK, dto.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
	return &authv1.RefreshResponse{
		Tokens: &authv1.TokenPair{
			AccessToken:      tokens.AccessToken,
			AccessExpiresAt:  tokens.AccessExpiresAt,
			RefreshToken:     tokens.RefreshToken,
			RefreshExpiresAt: tokens.RefreshExpiresAt,
			SessionId:        tokens.SessionId,
		},
	}, nil
}

func (h *Handler) Logout(ctx context.Context, r *authv1.LogoutRequest) (*emptypb.Empty, error) {
	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

	// var req dto.LogoutRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	if err := h.Auth.Logout(ctx, r.RefreshToken); err != nil {
		// status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			// h.respondError(w, status, "invalid or expired token", nil)
			return &emptypb.Empty{}, errors.New("invalid or expired token")
		}
		// h.respondError(w, status, err.Error(), nil)
		return &emptypb.Empty{}, err
	}

	// h.respond(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
	return &emptypb.Empty{}, nil
}

func (h *Handler) HandleLogoutAll(ctx context.Context, r *authv1.LogoutAllRequest) error {
	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

	// var req dto.LogoutRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	if err := h.Auth.LogoutAll(ctx, r.RefreshToken); err != nil {
		// status := http.StatusUnauthorized
		if errors.Is(err, app.ErrInvalidToken) {
			// h.respondError(w, status, "invalid or expired token", nil)
			return errors.New("invalid or expired token")
		}
		// h.respondError(w, status, err.Error(), nil)
		return err
	}

	// h.respond(w, http.StatusOK, dto.SuccessResponse{Status: "ok"})
	return nil
}

func (h *Handler) ValidateAccess(ctx context.Context, r *authv1.ValidateAccessRequest) (*authv1.ValidateAccessResponse, error) {
	claims, err := h.Auth.ValidateAccessToken(ctx, r.GetAccessToken())
	if err != nil {
		if errors.Is(err, app.ErrInvalidToken) {
			return nil, errors.New("invalid or expired token")
		}
		return nil, err
	}

	return &authv1.ValidateAccessResponse{
		UserId:       claims.Subject,
		SessionId:    claims.SessionID,
		TokenVersion: claims.TokenVersion,
	}, nil
}

// func (h *Handler) handleDomainError(w http.ResponseWriter, err error) {
// 	var (
// 		status  = http.StatusInternalServerError
// 		message = "internal server error"
// 		details []dto.ErrorDetail
// 	)

// 	switch {
// 	case errors.Is(err, context.Canceled):
// 		status = statusClientClosedRequest
// 		message = "request canceled"
// 	case errors.Is(err, context.DeadlineExceeded):
// 		status = http.StatusGatewayTimeout
// 		message = "request deadline exceeded"
// 	}

// 	var validationErrs *domain.ValidationErrors
// 	switch {
// 	case errors.As(err, &validationErrs):
// 		status = http.StatusBadRequest
// 		message = "validation error"
// 		details = dto.FromValidationErrors(validationErrs)
// 	case errors.Is(err, domain.ErrClientAlreadyExists):
// 		status = http.StatusConflict
// 		message = "client already exists"
// 	case errors.Is(err, domain.ErrClientNotFound):
// 		status = http.StatusNotFound
// 		message = "client not found"
// 	case errors.Is(err, domain.ErrAddressNotFound):
// 		status = http.StatusNotFound
// 		message = "address not found"
// 	default:
// 		message = err.Error()
// 	}

// 	h.logger.Warn("request failed", "status", status, "message", message, "error", err)
// 	h.respondError(w, status, message, details)
// }

// func (h *Handler) respondError(status int, message string, details []dto.ErrorDetail) {
// 	h.respond(status, dto.ErrorResponse{Error: message, Details: details})
// }

// func (h *Handler) respond(status int, payload any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	if payload == nil {
// 		return
// 	}
// 	_ = json.NewEncoder(w).Encode(payload)
// }
