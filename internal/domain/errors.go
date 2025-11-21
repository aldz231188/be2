package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrAddressNotFound     = errors.New("address not found")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if strings.TrimSpace(e.Field) == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors struct {
	errors []ValidationError
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{errors: make([]ValidationError, 0)}
}

func (ve *ValidationErrors) Add(field, message string) {
	if ve == nil {
		return
	}
	ve.errors = append(ve.errors, ValidationError{Field: field, Message: message})
}

func (ve *ValidationErrors) Merge(other *ValidationErrors) {
	if ve == nil || other == nil {
		return
	}
	ve.errors = append(ve.errors, other.errors...)
}

func (ve *ValidationErrors) HasErrors() bool {
	return ve != nil && len(ve.errors) > 0
}

func (ve *ValidationErrors) Error() string {
	if ve == nil || len(ve.errors) == 0 {
		return ""
	}
	parts := make([]string, len(ve.errors))
	for i, err := range ve.errors {
		parts[i] = err.Error()
	}
	return strings.Join(parts, ", ")
}

func (ve *ValidationErrors) Errors() []ValidationError {
	if ve == nil {
		return nil
	}
	return ve.errors
}
