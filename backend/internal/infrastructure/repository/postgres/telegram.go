package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type telegramRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewTelegramRepo(p *pg.Postgres) repository.TelegramRepository {
	return &telegramRepo{db: p.DB, builder: p.Builder}
}

func scanTelegramConn(row pgx.Row) (*entity.TelegramConnection, error) {
	c := &entity.TelegramConnection{}
	err := row.Scan(
		&c.ID, &c.UserID, &c.TelegramID, &c.ChatID, &c.Username,
		&c.VerificationCode, &c.VerifiedAt, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *telegramRepo) Create(ctx context.Context, conn *entity.TelegramConnection) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO telegram_connections(id, user_id, verification_code, created_at)
		VALUES($1, $2, $3, $4)
		ON CONFLICT(user_id) DO UPDATE SET verification_code=EXCLUDED.verification_code
	`, conn.ID, conn.UserID, conn.VerificationCode, conn.CreatedAt)
	if err != nil {
		return fmt.Errorf("telegramRepo.Create: %w", err)
	}
	return nil
}

func (r *telegramRepo) GetByUserID(ctx context.Context, userID string) (*entity.TelegramConnection, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, telegram_id, chat_id, username, verification_code, verified_at, created_at
		FROM telegram_connections WHERE user_id=$1
	`, userID)
	c, err := scanTelegramConn(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("telegram connection")
	}
	return c, err
}

func (r *telegramRepo) GetByChatID(ctx context.Context, chatID int64) (*entity.TelegramConnection, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, telegram_id, chat_id, username, verification_code, verified_at, created_at
		FROM telegram_connections WHERE chat_id=$1
	`, chatID)
	c, err := scanTelegramConn(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("telegram connection")
	}
	return c, err
}

func (r *telegramRepo) GetByVerificationCode(ctx context.Context, code string) (*entity.TelegramConnection, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, telegram_id, chat_id, username, verification_code, verified_at, created_at
		FROM telegram_connections WHERE verification_code=$1
	`, code)
	c, err := scanTelegramConn(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("telegram connection")
	}
	return c, err
}

func (r *telegramRepo) UpdateVerified(ctx context.Context, id string, telegramID, chatID int64, username string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE telegram_connections
		SET telegram_id=$2, chat_id=$3, username=$4, verified_at=NOW(), verification_code=NULL
		WHERE id=$1
	`, id, telegramID, chatID, username)
	return err
}

// Link upserts a verified connection: bot sends code → user enters on website.
func (r *telegramRepo) Link(ctx context.Context, userID string, telegramID, chatID int64, username string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO telegram_connections(id, user_id, telegram_id, chat_id, username, verified_at, created_at)
		VALUES(gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT(user_id) DO UPDATE
		SET telegram_id=$2, chat_id=$3, username=$4, verified_at=NOW(), verification_code=NULL
	`, userID, telegramID, chatID, username)
	return err
}

func (r *telegramRepo) Delete(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM telegram_connections WHERE user_id=$1`, userID)
	return err
}
