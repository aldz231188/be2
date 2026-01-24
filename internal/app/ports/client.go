package ports

import "context"

type ClientService interface {
	Create(ctx context.Context, userid, name, surename string) (string, error)
	// Delete(ctx context.Context, id string) error
}
