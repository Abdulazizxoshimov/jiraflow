package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type inviteRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewInviteRepo(p *pg.Postgres) repository.InviteRepository {
	return &inviteRepo{db: p.DB, builder: p.Builder}
}

const inviteCols = "id, email, role, token_hash, invited_by, expires_at, accepted_at, accepted_by, created_at"

func scanInvite(row pgx.Row) (*entity.Invite, error) {
	iv := &entity.Invite{}
	err := row.Scan(
		&iv.ID, &iv.Email, &iv.Role, &iv.TokenHash,
		&iv.InvitedBy, &iv.ExpiresAt, &iv.AcceptedAt, &iv.AcceptedBy, &iv.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

func (r *inviteRepo) Create(ctx context.Context, invite *entity.Invite) error {
	sql, args, err := r.builder.
		Insert("invites").
		Columns("id", "email", "role", "token_hash", "invited_by", "expires_at", "created_at").
		Values(invite.ID, invite.Email, invite.Role, invite.TokenHash, invite.InvitedBy, invite.ExpiresAt, invite.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("inviteRepo.Create: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("invite token already exists")
		}
		return fmt.Errorf("inviteRepo.Create: %w", err)
	}
	return nil
}

func (r *inviteRepo) GetByID(ctx context.Context, id string) (*entity.Invite, error) {
	sql, args, err := r.builder.
		Select(inviteCols).From("invites").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("inviteRepo.GetByID: %w", err)
	}
	iv, err := scanInvite(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("invite")
	}
	return iv, err
}

func (r *inviteRepo) GetByTokenHash(ctx context.Context, hash string) (*entity.Invite, error) {
	sql, args, err := r.builder.
		Select(inviteCols).From("invites").
		Where(sq.And{
			sq.Eq{"token_hash": hash},
			sq.Eq{"accepted_at": nil},
			sq.Expr("expires_at > NOW()"),
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("inviteRepo.GetByTokenHash: %w", err)
	}
	iv, err := scanInvite(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("invite")
	}
	return iv, err
}

func (r *inviteRepo) GetPendingByEmail(ctx context.Context, email string) (*entity.Invite, error) {
	sql, args, err := r.builder.
		Select(inviteCols).From("invites").
		Where(sq.And{
			sq.Eq{"email": email},
			sq.Eq{"accepted_at": nil},
			sq.Expr("expires_at > NOW()"),
		}).
		OrderBy("created_at DESC").Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("inviteRepo.GetPendingByEmail: %w", err)
	}
	iv, err := scanInvite(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("invite")
	}
	return iv, err
}

func (r *inviteRepo) List(ctx context.Context, filter *entity.Filter) ([]*entity.Invite, int, error) {
	where := sq.And{}
	if filter.Search != "" {
		where = append(where, sq.ILike{"email": "%" + filter.Search + "%"})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("invites").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("inviteRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(inviteCols).From("invites").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).
		ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("inviteRepo.List query: %w", err)
	}
	defer rows.Close()

	var invites []*entity.Invite
	for rows.Next() {
		iv, err := scanInvite(rows)
		if err != nil {
			return nil, 0, err
		}
		invites = append(invites, iv)
	}
	return invites, total, rows.Err()
}

func (r *inviteRepo) MarkAccepted(ctx context.Context, id, acceptedByUserID string) error {
	sql, args, err := r.builder.
		Update("invites").
		Set("accepted_at", sq.Expr("NOW()")).
		Set("accepted_by", acceptedByUserID).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("inviteRepo.MarkAccepted: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *inviteRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("invites").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("inviteRepo.Delete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *inviteRepo) DeleteExpired(ctx context.Context) error {
	sql, args, err := r.builder.
		Delete("invites").
		Where(sq.And{sq.Expr("expires_at < NOW()"), sq.Eq{"accepted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("inviteRepo.DeleteExpired: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
