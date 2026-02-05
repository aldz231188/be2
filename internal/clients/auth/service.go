package auth

import (
	"context"
	"errors"

	authv1 "be2/contracts/gen/auth/v1"
	"be2/internal/app/ports"
)

// type TokenPair struct {
// 	AccessToken      string
// 	AccessExpiresAt  int64
// 	RefreshToken     string
// 	RefreshExpiresAt int64
// 	SessionId        string
// }

type Service struct {
	c authv1.AuthServiceClient
}

func NewService(conn Conn) ports.AuthService {
	return &Service{c: authv1.NewAuthServiceClient(conn.ClientConn)}
}

func (s *Service) Register(ctx context.Context, login, password string) (*ports.TokenPair, error) {
	resp, err := s.c.Register(ctx, &authv1.RegisterRequest{Login: login, Password: password})
	if err != nil {
		return nil, err
	}
	a := resp.GetTokens()
	if a == nil {
		return nil, errors.New("authsvc returned empty token")
	}
	return &ports.TokenPair{
		AccessToken:      a.AccessToken,
		AccessExpiresAt:  a.AccessExpiresAt,
		RefreshToken:     a.RefreshToken,
		RefreshExpiresAt: a.RefreshExpiresAt,
		SessionId:        a.SessionId,
	}, nil
}
