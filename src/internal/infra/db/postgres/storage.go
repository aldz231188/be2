package db

import (
	"be2/internal/domain"
	sqlc "be2/internal/infra/db/sqlc_generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	q *sqlc.Queries
}

func NewUserRepo(pool *pgxpool.Pool) domain.Repo {
	return &UserRepo{q: sqlc.New(pool)}
}
