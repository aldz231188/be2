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

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}

func NewListener(cfg config.Config) (net.Listener, error) {
	return net.Listen("tcp", cfg.ClientSvcAddr)
}

func RegisterHandlers(s *grpc.Server, h *handlers.Handler) {
	clientv1.RegisterClientServiceServer(s, h)
}

func Run(lc fx.Lifecycle, s *grpc.Server, lis net.Listener, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("[clientsvc] gRPC on %s", cfg.ClientSvcAddr) // добавить проверку на пустой адрес
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
			return nil
		},
	})
}
