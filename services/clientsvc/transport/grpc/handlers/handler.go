package handlers

import (
	"be2/services/clientsvc/app"
	// "be2/services/clientsvc/domain"
	"be2/services/clientsvc/transport/grpc/mapper"
	"context"
	"encoding/json"
	// "errors"
	"log/slog"
	"net/http"

	clientv1 "be2/contracts/gen/client/v1"
	// "example.com/microdemo/clientsvc/services/clientsvc/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/protobuf/types/known/emptypb"
)

const statusClientClosedRequest = 499

type Handler struct {
	clientv1.UnimplementedClientServiceServer
	// AS     app.AddressService
	CS     app.ClientService
	logger *slog.Logger
}

func NewHandler(csi app.ClientService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		// AS:     asi,
		CS: csi,
		// Auth:   auth,
		logger: logger,
	}
}

func (h *Handler) CreateClient(ctx context.Context, req *clientv1.CreateClientRequest) (*clientv1.CreateClientResponse, error) {
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	// var clientRow mapper.CreateClientRequest
	// if err := json.NewDecoder(r.Body).Decode(&clientRow); err != nil {
	// 	h.respondError(w, http.StatusBadRequest, "invalid request body", nil)
	// 	return
	// }

	var clientRow mapper.CreateClientRequest
	clientRow.P = req
	client, err := clientRow.ToDomainClient()
	if err != nil {
		// h.handleDomainError(w, err)
		// return
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientID, err := h.CS.CreateClient(ctx, client)
	if err != nil {
		// h.handleDomainError(w, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	h.logger.InfoContext(ctx, "client created", "client_id", clientID)
	return &clientv1.CreateClientResponse{
		Clientid: clientID,
	}, nil

}

// func (h *Handler) HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
// 	defer cancel()

// 	var clientID mapper.UUIDRequest
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
// 		h.respondJSON(w, http.StatusOK, mapper.SuccessResponse{Status: "ok"})
// 	}
// }

// func (h *Handler) handleDomainError(w http.ResponseWriter, err error) {
// 	var (
// 		status  = http.StatusInternalServerError
// 		message = "internal server error"
// 		// details []mapper.ErrorDetail
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
// 		details = mapper.FromValidationErrors(validationErrs)
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
// 	// h.respondError(w, status, message, details)
// }

func (h *Handler) respondError(w http.ResponseWriter, status int, message string, details []mapper.ErrorDetail) {
	h.respondJSON(w, status, mapper.ErrorResponse{Error: message, Details: details})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}
