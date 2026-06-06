package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type auditRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewAuditRepo(p *pg.Postgres) repository.AuditRepository {
	return &auditRepo{db: p.DB, builder: p.Builder}
}

func (r *auditRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	detailsJSON, err := json.Marshal(log.Details)
	if err != nil {
		return fmt.Errorf("auditRepo.Create marshal details: %w", err)
	}
	sql, args, err := r.builder.
		Insert("audit_logs").
		Columns("id", "user_id", "action", "entity_type", "entity_id", "details", "ip_address", "user_agent", "created_at").
		Values(log.ID, log.UserID, log.Action, log.EntityType, log.EntityID, detailsJSON, log.IPAddress, log.UserAgent, log.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("auditRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *auditRepo) List(ctx context.Context, filter *entity.AuditLogFilter) ([]*entity.AuditLog, int, error) {
	where := sq.And{}
	if filter.UserID != "" {
		where = append(where, sq.Eq{"user_id": filter.UserID})
	}
	if filter.Action != "" {
		where = append(where, sq.Eq{"action": filter.Action})
	}
	if filter.EntityType != "" {
		where = append(where, sq.Eq{"entity_type": filter.EntityType})
	}
	if filter.EntityID != "" {
		where = append(where, sq.Eq{"entity_id": filter.EntityID})
	}
	if filter.CreatedFrom != nil {
		where = append(where, sq.GtOrEq{"created_at": filter.CreatedFrom})
	}
	if filter.CreatedTo != nil {
		where = append(where, sq.LtOrEq{"created_at": filter.CreatedTo})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("audit_logs").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("auditRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select("id", "user_id", "action", "entity_type", "entity_id", "details", "ip_address", "user_agent", "created_at").
		From("audit_logs").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("auditRepo.List query: %w", err)
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		l := &entity.AuditLog{}
		var detailsJSON []byte
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID,
			&detailsJSON, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(detailsJSON, &l.Details)
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}
