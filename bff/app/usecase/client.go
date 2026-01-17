package usecase

import (
	"context"
	// "strings"
	"time"

	"be2/internal/app/ports"
	"be2/internal/config"
	// "be2/internal/domain"
)

type ClientUsecase struct {
	svc        ports.ClientService
	depTimeout time.Duration
}

func NewClientUsecase(cfg config.Config, svc ports.ClientService) ClientUsecase {
	return ClientUsecase{svc: svc, depTimeout: cfg.DepTimeout}
}

func (u *ClientUsecase) Create(ctx context.Context, userid, name, surename string) (string, error) {
	// client = strings.TrimSpace(client)
	// if client == "" {
	// 	return "", domain.ValidationError{} //другая ошибка
	// }

	cctx, cancel := context.WithTimeout(ctx, u.depTimeout)
	defer cancel()

	return u.svc.Create(cctx, userid, name, surename)
}

// func (u *ClientUsecase) Delete(ctx context.Context, id string) error {
// 	id = strings.TrimSpace(id)
// 	if id == "" {
// 		return domain.ValidationError{}
// 	}

// 	cctx, cancel := context.WithTimeout(ctx, u.depTimeout)
// 	defer cancel()

// 	return u.svc.Delete(cctx, id)
// }
