package repo

import (
	"be2/services/auth/internal/domain"
	store "be2/services/auth/internal/infra/db/sqlc_generated"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repo) CreateSession(ctx context.Context, session domain.Session) error {
	err := r.q.CreateSession(ctx, r.toSessionParams(session))
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *Repo) GetSessionByHash(ctx context.Context, hash string) (domain.Session, error) {
	row, err := r.q.GetSessionByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return domain.Session{}, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, domain.ErrSessionNotFound
		}
		return domain.Session{}, fmt.Errorf("get session: %w", err)
	}

	return domain.Session{
		JTIHash:   row.JtiHash,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt.Time,
		RevokedAt: nullableTime(row.RevokedAt),
		CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *Repo) RevokeSession(ctx context.Context, hash string) error {
	_, err := r.q.RevokeSession(ctx, hash)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (r *Repo) RevokeSessionsByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.q.RevokeSessionsByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("revoke user sessions: %w", err)
	}
	return nil
}

func (r *Repo) toSessionParams(session domain.Session) store.CreateSessionParams {
	return store.CreateSessionParams{
		JtiHash: session.JTIHash,
		UserID:  session.UserID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  session.ExpiresAt,
			Valid: true,
		},
	}
}

func nullableTime(ts pgtype.Timestamptz) *time.Time {
	if ts.Valid {
		return &ts.Time
	}
	return nil
}
