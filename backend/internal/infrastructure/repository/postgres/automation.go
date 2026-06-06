package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type automationRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewAutomationRepo(p *pg.Postgres) repository.AutomationRepository {
	return &automationRepo{db: p.DB, builder: p.Builder}
}

func (r *automationRepo) Create(ctx context.Context, rule *entity.AutomationRule) error {
	condJSON, _ := json.Marshal(rule.Conditions)
	actJSON, _ := json.Marshal(rule.Actions)
	cfgJSON, _ := json.Marshal(rule.TriggerConfig)

	sql, args, err := r.builder.
		Insert("automation_rules").
		Columns("id", "project_id", "name", "description", "trigger_type", "trigger_config",
			"conditions", "actions", "is_active", "created_by", "created_at", "updated_at").
		Values(rule.ID, rule.ProjectID, rule.Name, rule.Description, rule.TriggerType, cfgJSON,
			condJSON, actJSON, rule.IsActive, rule.CreatedBy, rule.CreatedAt, rule.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("automationRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *automationRepo) GetByID(ctx context.Context, id string) (*entity.AutomationRule, error) {
	sql, args, err := r.builder.
		Select("id", "project_id", "name", "description", "trigger_type", "trigger_config",
			"conditions", "actions", "is_active", "created_by", "created_at", "updated_at").
		From("automation_rules").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("automationRepo.GetByID build: %w", err)
	}
	row := r.db.QueryRow(ctx, sql, args...)
	rule, err := scanAutomationRule(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("automation rule")
	}
	return rule, err
}

func (r *automationRepo) List(ctx context.Context, filter *entity.AutomationFilter) ([]*entity.AutomationRule, int, error) {
	q := r.builder.
		Select("id", "project_id", "name", "description", "trigger_type", "trigger_config",
			"conditions", "actions", "is_active", "created_by", "created_at", "updated_at").
		From("automation_rules").
		OrderBy("created_at DESC")

	if filter.ProjectID != "" {
		q = q.Where(sq.Eq{"project_id": filter.ProjectID})
	}
	if filter.TriggerType != "" {
		q = q.Where(sq.Eq{"trigger_type": filter.TriggerType})
	}
	if filter.IsActive != nil {
		q = q.Where(sq.Eq{"is_active": *filter.IsActive})
	}

	cntSQL, cntArgs, _ := q.Columns().ToSql()
	_ = cntSQL
	_ = cntArgs

	countQ := r.builder.Select("COUNT(*)").From("automation_rules")
	if filter.ProjectID != "" {
		countQ = countQ.Where(sq.Eq{"project_id": filter.ProjectID})
	}
	if filter.IsActive != nil {
		countQ = countQ.Where(sq.Eq{"is_active": *filter.IsActive})
	}
	countSQL, countArgs, _ := countQ.ToSql()
	var total int
	_ = r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)

	limit := filter.GetLimit()
	offset := filter.Offset()
	q = q.Limit(uint64(limit)).Offset(uint64(offset))

	listSQL, listArgs, err := q.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("automationRepo.List build: %w", err)
	}

	rows, err := r.db.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []*entity.AutomationRule
	for rows.Next() {
		rule, err := scanAutomationRule(rows)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, rule)
	}
	return rules, total, nil
}

func (r *automationRepo) Update(ctx context.Context, rule *entity.AutomationRule) error {
	condJSON, _ := json.Marshal(rule.Conditions)
	actJSON, _ := json.Marshal(rule.Actions)
	cfgJSON, _ := json.Marshal(rule.TriggerConfig)

	sql, args, err := r.builder.
		Update("automation_rules").
		Set("name", rule.Name).
		Set("description", rule.Description).
		Set("trigger_type", rule.TriggerType).
		Set("trigger_config", cfgJSON).
		Set("conditions", condJSON).
		Set("actions", actJSON).
		Set("is_active", rule.IsActive).
		Set("updated_at", rule.UpdatedAt).
		Where(sq.Eq{"id": rule.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("automationRepo.Update build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *automationRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("automation_rules").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("automationRepo.Delete build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *automationRepo) SetActive(ctx context.Context, id string, isActive bool) error {
	sql, args, err := r.builder.
		Update("automation_rules").
		Set("is_active", isActive).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("automationRepo.SetActive build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *automationRepo) FindByTrigger(ctx context.Context, projectID, triggerType string) ([]*entity.AutomationRule, error) {
	sql, args, err := r.builder.
		Select("id", "project_id", "name", "description", "trigger_type", "trigger_config",
			"conditions", "actions", "is_active", "created_by", "created_at", "updated_at").
		From("automation_rules").
		Where(sq.And{
			sq.Eq{"project_id": projectID},
			sq.Eq{"trigger_type": triggerType},
			sq.Eq{"is_active": true},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("automationRepo.FindByTrigger build: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*entity.AutomationRule
	for rows.Next() {
		rule, err := scanAutomationRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *automationRepo) SaveLog(ctx context.Context, log *entity.AutomationLog) error {
	sql, args, err := r.builder.
		Insert("automation_logs").
		Columns("id", "rule_id", "entity_id", "entity_type", "status", "executed_at", "error_msg").
		Values(log.ID, log.RuleID, log.EntityID, log.EntityType, log.Status, log.ExecutedAt, log.ErrorMsg).
		ToSql()
	if err != nil {
		return fmt.Errorf("automationRepo.SaveLog build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *automationRepo) ListLogs(ctx context.Context, ruleID string, filter *entity.Filter) ([]*entity.AutomationLog, int, error) {
	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("automation_logs").Where(sq.Eq{"rule_id": ruleID}).ToSql()
	_ = r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total)

	listSQL, listArgs, err := r.builder.
		Select("id", "rule_id", "entity_id", "entity_type", "status", "executed_at", "error_msg").
		From("automation_logs").
		Where(sq.Eq{"rule_id": ruleID}).
		OrderBy("executed_at DESC").
		Limit(uint64(filter.GetLimit())).
		Offset(uint64(filter.Offset())).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("automationRepo.ListLogs build: %w", err)
	}

	rows, err := r.db.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*entity.AutomationLog
	for rows.Next() {
		l := &entity.AutomationLog{}
		if err := rows.Scan(&l.ID, &l.RuleID, &l.EntityID, &l.EntityType, &l.Status, &l.ExecutedAt, &l.ErrorMsg); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

type automationScanner interface {
	Scan(dest ...any) error
}

func scanAutomationRule(row automationScanner) (*entity.AutomationRule, error) {
	r := &entity.AutomationRule{}
	var condJSON, actJSON, cfgJSON []byte
	err := row.Scan(
		&r.ID, &r.ProjectID, &r.Name, &r.Description, &r.TriggerType, &cfgJSON,
		&condJSON, &actJSON, &r.IsActive, &r.CreatedBy, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(cfgJSON) > 0 {
		_ = json.Unmarshal(cfgJSON, &r.TriggerConfig)
	}
	if len(condJSON) > 0 {
		_ = json.Unmarshal(condJSON, &r.Conditions)
	}
	if len(actJSON) > 0 {
		_ = json.Unmarshal(actJSON, &r.Actions)
	}
	if r.TriggerConfig == nil {
		r.TriggerConfig = map[string]any{}
	}
	return r, nil
}
