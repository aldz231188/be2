package authserver

import (
	"be2/internal/app"
	"be2/internal/grpc/authpb"
	"context"
	"errors"
	"log/slog"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	service app.AuthService
}

func NewServer(service app.AuthService) *Server {
	return &Server{service: service}
}

func Register(lc fx.Lifecycle, cfg Config, server *Server, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, server)

	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("auth grpc server listening", "addr", cfg.Addr)
			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					logger.Error("auth grpc server stopped", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			stopped := make(chan struct{})
			go func() {
				grpcServer.GracefulStop()
				close(stopped)
			}()
			select {
			case <-ctx.Done():
				grpcServer.Stop()
				return ctx.Err()
			case <-stopped:
				return nil
			}
		},
	})

	return nil
}

func (s *Server) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.TokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	pair, err := s.service.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, app.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.TokenPairResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.TokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	pair, err := s.service.Register(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, app.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		case errors.Is(err, app.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &authpb.TokenPairResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.TokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	pair, err := s.service.Refresh(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, app.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.TokenPairResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) LogoutCurrent(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	if err := s.service.LogoutCurrent(ctx, req.RefreshToken); err != nil {
		if errors.Is(err, app.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.LogoutResponse{Status: "ok"}, nil
}

func (s *Server) LogoutAll(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	if err := s.service.LogoutAll(ctx, req.RefreshToken); err != nil {
		if errors.Is(err, app.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.LogoutResponse{Status: "ok"}, nil
}

func (s *Server) ValidateAccessToken(ctx context.Context, req *authpb.ValidateAccessTokenRequest) (*authpb.ValidateAccessTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}

	claims, err := s.service.ValidateAccessToken(ctx, req.AccessToken)
	if err != nil {
		if errors.Is(err, app.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.ValidateAccessTokenResponse{Claims: toProtoClaims(claims)}, nil
}

func toProtoClaims(claims *app.TokenClaims) *authpb.TokenClaims {
	if claims == nil {
		return nil
	}

	protoClaims := &authpb.TokenClaims{
		TokenType:    claims.TokenType,
		TokenVersion: claims.TokenVersion,
		SessionId:    claims.SessionID,
		Subject:      claims.Subject,
		Id:           claims.ID,
	}

	if claims.ExpiresAt != nil {
		protoClaims.ExpiresAtUnix = claims.ExpiresAt.Time.Unix()
	}
	if claims.IssuedAt != nil {
		protoClaims.IssuedAtUnix = claims.IssuedAt.Time.Unix()
	}

	return protoClaims
}

var _ authpb.AuthServiceServer = (*Server)(nil)
