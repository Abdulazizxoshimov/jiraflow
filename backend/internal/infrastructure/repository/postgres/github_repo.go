package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type githubRepoRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewGitHubRepoRepo(p *pg.Postgres) repository.GitHubRepoRepository {
	return &githubRepoRepo{db: p.DB, builder: p.Builder}
}

func scanGitHubRepo(row pgx.Row) (*entity.GitHubRepo, error) {
	r := &entity.GitHubRepo{}
	err := row.Scan(&r.ID, &r.ProjectID, &r.RepoFullName, &r.RepoURL, &r.WebhookSecret, &r.CreatedBy, &r.CreatedAt)
	return r, err
}

func (r *githubRepoRepo) Create(ctx context.Context, repo *entity.GitHubRepo) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO github_repos(id, project_id, repo_full_name, repo_url, webhook_secret, created_by, created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
	`, repo.ID, repo.ProjectID, repo.RepoFullName, repo.RepoURL, repo.WebhookSecret, repo.CreatedBy, repo.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("github repo already connected to this project")
		}
		return fmt.Errorf("githubRepoRepo.Create: %w", err)
	}
	return nil
}

func (r *githubRepoRepo) GetByProjectID(ctx context.Context, projectID string) (*entity.GitHubRepo, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, project_id, repo_full_name, repo_url, webhook_secret, created_by, created_at FROM github_repos WHERE project_id=$1`,
		projectID)
	repo, err := scanGitHubRepo(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("github repo")
	}
	return repo, err
}

func (r *githubRepoRepo) GetByRepoFullName(ctx context.Context, repoFullName string) (*entity.GitHubRepo, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, project_id, repo_full_name, repo_url, webhook_secret, created_by, created_at FROM github_repos WHERE repo_full_name=$1`,
		repoFullName)
	repo, err := scanGitHubRepo(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("github repo")
	}
	return repo, err
}

func (r *githubRepoRepo) Delete(ctx context.Context, projectID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM github_repos WHERE project_id=$1`, projectID)
	return err
}
