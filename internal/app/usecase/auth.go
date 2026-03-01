package usecase

import (
	"context"
	// "strings"
	"time"

	"be2/internal/app/ports"
	"be2/internal/config"
	// "be2/internal/domain"
)

type AuthUsecase struct {
	svc        ports.AuthService
	depTimeout time.Duration
}

func NewAuthUsecase(cfg config.Config, svc ports.AuthService) AuthUsecase {
	return AuthUsecase{svc: svc, depTimeout: cfg.DepTimeout}
}

func (u *AuthUsecase) Register(ctx context.Context, login, password string) (*ports.TokenPair, error) {
	// Auth = strings.TrimSpace(Auth)
	// if Auth == "" {
	// 	return "", domain.ValidationError{} //другая ошибка
	// }

	cctx, cancel := context.WithTimeout(ctx, u.depTimeout)
	defer cancel()

	return u.svc.Register(cctx, login, password)
}
func (u *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	// Auth = strings.TrimSpace(Auth)
	// if Auth == "" {
	// 	return "", domain.ValidationError{} //другая ошибка
	// }

	cctx, cancel := context.WithTimeout(ctx, u.depTimeout)
	defer cancel()

	return u.svc.Logout(cctx, refreshToken)
}

// func (u *AuthUsecase) Delete(ctx context.Context, id string) error {
// 	id = strings.TrimSpace(id)
// 	if id == "" {
// 		return domain.ValidationError{}
// 	}

// 	cctx, cancel := context.WithTimeout(ctx, u.depTimeout)
// 	defer cancel()

// 	return u.svc.Delete(cctx, id)
// }
