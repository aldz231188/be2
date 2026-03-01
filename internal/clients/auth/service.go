package auth

import (
	"context"
	"errors"

	authv1 "be2/contracts/gen/auth/v1"
	"be2/internal/app/ports"
)

type Service struct {
	c authv1.AuthServiceClient
}

func NewService(conn Conn) ports.AuthService {
	return &Service{c: authv1.NewAuthServiceClient(conn.ClientConn)}
}

var _ ports.AuthService = (*Service)(nil)

func (s *Service) Register(ctx context.Context, login, password string) (*ports.TokenPair, error) {
	resp, err := s.c.Register(ctx, &authv1.RegisterRequest{Login: login, Password: password})
	if err != nil {
		return nil, err
	}
	tokens := resp.GetTokens()
	if tokens == nil {
		return nil, errors.New("authsvc returned empty token")
	}
	return &ports.TokenPair{
		AccessToken:      tokens.AccessToken,
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshToken:     tokens.RefreshToken,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
		SessionId:        tokens.SessionId,
	}, nil
}
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.c.Logout(ctx, &authv1.LogoutRequest{RefreshToken: refreshToken})
	if err != nil {
		return err
	}
	return nil

}

func (s *Service) ValidateAccess(ctx context.Context, accessToken string) (string, error) {
	resp, err := s.c.ValidateAccess(ctx, &authv1.ValidateAccessRequest{AccessToken: accessToken})
	if err != nil {
		return "", err
	}
	if resp.GetUserId() == "" {
		return "", errors.New("authsvc returned empty user id")
	}
	return resp.GetUserId(), nil
}
