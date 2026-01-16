package client

import (
	"context"
	"errors"

	clientv1 "be2/contracts/gen/client/v1"
	"be2/internal/app/ports"
)

type Service struct {
	c clientv1.ClientServiceClient
}

func NewService(conn Conn) ports.ClientService {
	return &Service{c: clientv1.NewClientServiceClient(conn.ClientConn)}
}

func (s *Service) Create(ctx context.Context, userid, name, surename string) (string, error) {
	resp, err := s.c.CreateClient(ctx, &clientv1.CreateClientRequest{Userid: userid, Name: name, Surname: surename})
	if err != nil {
		return "", err
	}
	u := resp.GetClientid()
	if u == "" {
		return "", errors.New("clientsvc returned empty client")
	}
	return u, nil
}

// func (s *Service) Delete(ctx context.Context, id string) error {
// 	_, err := s.c.DeleteClient(ctx, &clientv1.DeleteClientRequest{Id: id})
// 	return err
// }
