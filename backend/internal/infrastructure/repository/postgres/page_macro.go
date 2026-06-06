package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageMacroRepo struct {
	db *pgxpool.Pool
}

func NewPageMacroRepo(p *pg.Postgres) repository.PageMacroRepository {
	return &pageMacroRepo{db: p.DB}
}

func (r *pageMacroRepo) Upsert(ctx context.Context, macro *entity.PageMacro) error {
	if macro.ID == "" {
		macro.ID = uuid.NewString()
	}
	configJSON, err := json.Marshal(macro.Config)
	if err != nil {
		return fmt.Errorf("pageMacroRepo.Upsert marshal: %w", err)
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO page_macros(id, page_id, macro_type, config, created_at)
		VALUES($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE SET macro_type=EXCLUDED.macro_type, config=EXCLUDED.config
	`, macro.ID, macro.PageID, macro.MacroType, configJSON)
	return err
}

func (r *pageMacroRepo) ListByPage(ctx context.Context, pageID string) ([]*entity.PageMacro, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, page_id, macro_type, config, created_at
		FROM page_macros WHERE page_id=$1 ORDER BY created_at ASC
	`, pageID)
	if err != nil {
		return nil, fmt.Errorf("pageMacroRepo.ListByPage: %w", err)
	}
	defer rows.Close()

	var macros []*entity.PageMacro
	for rows.Next() {
		m := &entity.PageMacro{}
		var configJSON []byte
		if err := rows.Scan(&m.ID, &m.PageID, &m.MacroType, &configJSON, &m.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(configJSON, &m.Config)
		macros = append(macros, m)
	}
	return macros, rows.Err()
}

func (r *pageMacroRepo) GetByID(ctx context.Context, id string) (*entity.PageMacro, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, page_id, macro_type, config, created_at
		FROM page_macros WHERE id=$1
	`, id)
	m := &entity.PageMacro{}
	var configJSON []byte
	if err := row.Scan(&m.ID, &m.PageID, &m.MacroType, &configJSON, &m.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("page macro")
		}
		return nil, err
	}
	_ = json.Unmarshal(configJSON, &m.Config)
	return m, nil
}

func (r *pageMacroRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_macros WHERE id=$1`, id)
	return err
}
