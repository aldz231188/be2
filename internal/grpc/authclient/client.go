package authclient

import (
	"be2/internal/app"
	"be2/internal/domain"
	"be2/internal/grpc/authpb"
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type client struct {
	client  authpb.AuthServiceClient
	timeout time.Duration
}

func NewAuthClient(lc fx.Lifecycle, cfg Config, logger *slog.Logger) (app.AuthService, error) {
	if logger == nil {
		logger = slog.Default()
	}

	conn, err := grpc.Dial(cfg.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("auth grpc client connected", "target", cfg.Target)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return &client{client: authpb.NewAuthServiceClient(conn), timeout: cfg.Timeout}, nil
}

func (c *client) Authenticate(ctx context.Context, username, password string) (app.TokenPair, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Authenticate(ctx, &authpb.AuthenticateRequest{Username: username, Password: password})
	if err != nil {
		return app.TokenPair{}, mapCredentialsError(err)
	}

	return app.TokenPair{AccessToken: resp.AccessToken, RefreshToken: resp.RefreshToken}, nil
}

func (c *client) Register(ctx context.Context, username, password string) (app.TokenPair, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Register(ctx, &authpb.RegisterRequest{Username: username, Password: password})
	if err != nil {
		return app.TokenPair{}, mapRegisterError(err)
	}

	return app.TokenPair{AccessToken: resp.AccessToken, RefreshToken: resp.RefreshToken}, nil
}

func (c *client) Refresh(ctx context.Context, refreshToken string) (app.TokenPair, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Refresh(ctx, &authpb.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return app.TokenPair{}, mapTokenError(err)
	}

	return app.TokenPair{AccessToken: resp.AccessToken, RefreshToken: resp.RefreshToken}, nil
}

func (c *client) LogoutCurrent(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.LogoutCurrent(ctx, &authpb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil {
		return mapTokenError(err)
	}
	return nil
}

func (c *client) LogoutAll(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.client.LogoutAll(ctx, &authpb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil {
		return mapTokenError(err)
	}
	return nil
}

func (c *client) ValidateAccessToken(ctx context.Context, token string) (*app.TokenClaims, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.ValidateAccessToken(ctx, &authpb.ValidateAccessTokenRequest{AccessToken: token})
	if err != nil {
		return nil, mapTokenError(err)
	}
	if resp == nil || resp.Claims == nil {
		return nil, app.ErrInvalidToken
	}

	return toAppClaims(resp.Claims), nil
}

func mapCredentialsError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch s.Code() {
	case codes.Unauthenticated:
		return app.ErrInvalidCredentials
	case codes.InvalidArgument:
		return app.ErrInvalidCredentials
	default:
		return err
	}
}

func mapTokenError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch s.Code() {
	case codes.Unauthenticated, codes.PermissionDenied, codes.NotFound:
		return app.ErrInvalidToken
	default:
		return err
	}
}

func mapRegisterError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch s.Code() {
	case codes.AlreadyExists:
		return domain.ErrUserAlreadyExists
	case codes.InvalidArgument:
		return app.ErrInvalidCredentials
	case codes.Unauthenticated:
		return app.ErrInvalidCredentials
	default:
		return err
	}
}

func toAppClaims(claims *authpb.TokenClaims) *app.TokenClaims {
	if claims == nil {
		return nil
	}

	appClaims := &app.TokenClaims{
		TokenType:    claims.TokenType,
		TokenVersion: claims.TokenVersion,
		SessionID:    claims.SessionId,
	}

	appClaims.Subject = claims.Subject
	appClaims.ID = claims.Id
	if claims.ExpiresAtUnix != 0 {
		appClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(claims.ExpiresAtUnix, 0))
	}
	if claims.IssuedAtUnix != 0 {
		appClaims.IssuedAt = jwt.NewNumericDate(time.Unix(claims.IssuedAtUnix, 0))
	}

	return appClaims
}
