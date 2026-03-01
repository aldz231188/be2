package repo

import (
	"be2/services/clientsvc/domain"
	store "be2/services/clientsvc/infra/db/sqlc_generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ q *store.Queries }

var (
	_ domain.ClientRepo = (*Repo)(nil)

// _ domain.SupplierRepo = (*Repo)(nil)
)

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{q: store.New(pool)} }
