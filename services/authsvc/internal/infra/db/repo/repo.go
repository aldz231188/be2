package repo

import (
	"be2/services/authsvc/internal/domain"
	store "be2/services/authsvc/internal/infra/db/sqlc_generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{ q *store.Queries }

var (
	_ domain.UserRepo    = (*Repo)(nil)
	_ domain.SessionRepo = (*Repo)(nil)

// _ domain.SupplierRepo = (*Repo)(nil)
)

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{q: store.New(pool)} }
