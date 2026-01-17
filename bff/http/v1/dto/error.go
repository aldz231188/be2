package dto

import "be2/internal/domain"

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string        `json:"error"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status string `json:"status"`
}

func FromValidationErrors(errs *domain.ValidationErrors) []ErrorDetail {
	if errs == nil {
		return nil
	}
	details := make([]ErrorDetail, 0, len(errs.Errors()))
	for _, e := range errs.Errors() {
		details = append(details, ErrorDetail{
			Field:   e.Field,
			Message: e.Message,
		})
	}
	return details
}
