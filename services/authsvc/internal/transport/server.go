package server

import (
	"context"
	"log"
	"net"

	authv1 "be2/contracts/gen/auth/v1"
	"be2/services/authsvc/internal/config"
	"be2/services/authsvc/internal/transport/grpc/handlers"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}

func NewListener(cfg config.Config) (net.Listener, error) {
	return net.Listen("tcp", cfg.authSvcAddr)
}

func RegisterHandlers(s *grpc.Server, h *handlers.Handler) {
	authv1.RegisterauthServiceServer(s, h)
}

func Run(lc fx.Lifecycle, s *grpc.Server, lis net.Listener, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("[authsvc] gRPC on %s", cfg.authSvcAddr) // добавить проверку на пустой адрес
			go func() {
				if err := s.Serve(lis); err != nil {
					log.Printf("[authsvc] serve error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("[authsvc] stopping")
			s.GracefulStop()
			return nil
		},
	})
}
