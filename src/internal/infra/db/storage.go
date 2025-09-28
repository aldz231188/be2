package db

import (
	"be2/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	q *Queries
}

func NewUserRepo(pool *pgxpool.Pool) domain.Repo {
	return &UserRepo{q: New(pool)}
}
