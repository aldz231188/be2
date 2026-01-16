package server

import (
	"context"
	"log"
	"net"

	clientv1 "be2/contracts/gen/client/v1"
	"be2/services/clientsvc/config"
	"be2/services/clientsvc/transport/grpc/handlers"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Start(lc fx.Lifecycle, s *grpc.Server, lis net.Listener, h *handlers.Handler, cfg config.Config) {
	Register(s, h)       // важно: регистрируем до Serve
	Run(lc, s, lis, cfg) // вешаем lifecycle + Serve
}

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}

func NewListener(cfg config.Config) (net.Listener, error) {
	return net.Listen("tcp", cfg.ClientSvcAddr)
}

func Register(s *grpc.Server, h *handlers.Handler) {
	clientv1.RegisterClientServiceServer(s, h)
}

func Run(lc fx.Lifecycle, s *grpc.Server, lis net.Listener, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("[clientsvc] gRPC on %s", cfg.ClientSvcAddr)
			go func() {
				if err := s.Serve(lis); err != nil {
					log.Printf("[clientsvc] serve error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("[clientsvc] stopping")
			s.GracefulStop()
			_ = lis.Close()
			return nil
		},
	})
}
