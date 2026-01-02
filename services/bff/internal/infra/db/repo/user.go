package repo

import (
	"be2/services/bff/internal/domain"
	"be2/services/bff/internal/infra/db/sqlc_generated"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return domain.User{}, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("get user by username: %w", err)
	}

	return domain.User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		TokenVersion: row.TokenVersion,
	}, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return domain.User{}, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return domain.User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		TokenVersion: row.TokenVersion,
	}, nil
}

func (r *Repo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, sqlc_generated.CreateUserParams{Username: user.Username, PasswordHash: user.PasswordHash}) //некрасиво как то
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			return domain.User{}, err
		case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
			return domain.User{}, domain.ErrUserAlreadyExists
		default:
			return domain.User{}, fmt.Errorf("create user: %w", err)
		}
	}

	return domain.User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		TokenVersion: row.TokenVersion,
	}, nil
}

func (r *Repo) IncrementTokenVersion(ctx context.Context, id uuid.UUID) error {
	_, err := r.q.IncrementTokenVersion(ctx, id)
	if err != nil {
		return fmt.Errorf("increment token version: %w", err)
	}
	return nil
}
