package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageLockRepo struct {
	db *pgxpool.Pool
}

func NewPageLockRepo(p *pg.Postgres) repository.PageLockRepository {
	return &pageLockRepo{db: p.DB}
}

func (r *pageLockRepo) Acquire(ctx context.Context, lock *entity.PageLock) error {
	sessionID := lock.SessionID
	if sessionID == "" {
		sessionID = lock.UserID
	}
	// Allow take-over only when lock is expired OR the same user is re-acquiring.
	// This whole operation is atomic — no TOCTOU race between Get() and INSERT.
	tag, err := r.db.Exec(ctx, `
		INSERT INTO page_locks(page_id, user_id, session_id, expires_at, created_at)
		VALUES($1, $2, $3, $4, NOW())
		ON CONFLICT (page_id) DO UPDATE
		  SET user_id    = EXCLUDED.user_id,
		      session_id = EXCLUDED.session_id,
		      expires_at = EXCLUDED.expires_at,
		      created_at = NOW()
		  WHERE page_locks.expires_at < NOW()
		     OR page_locks.user_id = EXCLUDED.user_id
	`, lock.PageID, lock.UserID, sessionID, lock.ExpiresAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperr.Conflict("page is locked by another user")
	}
	return nil
}

func (r *pageLockRepo) Release(ctx context.Context, pageID, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_locks WHERE page_id=$1 AND user_id=$2`, pageID, userID)
	return err
}

func (r *pageLockRepo) Get(ctx context.Context, pageID string) (*entity.PageLock, error) {
	row := r.db.QueryRow(ctx, `
		SELECT page_id, user_id, expires_at, created_at
		FROM page_locks WHERE page_id=$1 AND expires_at > NOW()
	`, pageID)
	lock := &entity.PageLock{}
	if err := row.Scan(&lock.PageID, &lock.UserID, &lock.ExpiresAt, &lock.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("page lock")
		}
		return nil, err
	}
	return lock, nil
}

func (r *pageLockRepo) Extend(ctx context.Context, pageID, userID string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE page_locks SET expires_at=$1
		WHERE page_id=$2 AND user_id=$3 AND expires_at > NOW()
	`, expiresAt, pageID, userID)
	return err
}

func (r *pageLockRepo) Cleanup(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_locks WHERE expires_at < NOW()`)
	return err
}
