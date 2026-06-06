package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	ghinfra "github.com/jira-backend/jiraflow-backend/internal/infrastructura/github"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	repoRepo     repository.GitHubRepoRepository
	commitRepo   repository.IssueCommitRepository
	prRepo       repository.IssuePRRepository
	issueRepo    repository.IssueRepository
	projectRepo  repository.ProjectRepository
	workflowRepo repository.WorkflowRepository
}

func New(
	repoRepo repository.GitHubRepoRepository,
	commitRepo repository.IssueCommitRepository,
	prRepo repository.IssuePRRepository,
	issueRepo repository.IssueRepository,
	projectRepo repository.ProjectRepository,
	workflowRepo repository.WorkflowRepository,
) UseCase {
	return &useCase{
		repoRepo:     repoRepo,
		commitRepo:   commitRepo,
		prRepo:       prRepo,
		issueRepo:    issueRepo,
		projectRepo:  projectRepo,
		workflowRepo: workflowRepo,
	}
}

func (uc *useCase) ConnectRepo(ctx context.Context, projectID, userID string, req *entity.ConnectRepoReq) (*entity.GitHubRepo, error) {
	repo := &entity.GitHubRepo{
		ID:            uuid.NewString(),
		ProjectID:     projectID,
		RepoFullName:  req.RepoFullName,
		RepoURL:       req.RepoURL,
		WebhookSecret: req.WebhookSecret,
		CreatedBy:     userID,
		CreatedAt:     time.Now().UTC(),
	}
	if err := uc.repoRepo.Create(ctx, repo); err != nil {
		return nil, fmt.Errorf("github.ConnectRepo: %w", err)
	}
	return repo, nil
}

func (uc *useCase) DisconnectRepo(ctx context.Context, projectID string) error {
	return uc.repoRepo.Delete(ctx, projectID)
}

func (uc *useCase) GetRepo(ctx context.Context, projectID string) (*entity.GitHubRepo, error) {
	return uc.repoRepo.GetByProjectID(ctx, projectID)
}

func (uc *useCase) HandlePushEvent(ctx context.Context, repoFullName string, commits []ghinfra.Commit) error {
	ghRepo, err := uc.repoRepo.GetByRepoFullName(ctx, repoFullName)
	if err != nil {
		return nil // repo not connected — ignore
	}

	for _, c := range commits {
		keys := ghinfra.ExtractIssueKeys(c.Message)
		for _, key := range keys {
			issue, err := uc.issueRepo.GetByKey(ctx, key)
			if err != nil {
				continue
			}
			var committedAt *time.Time
			if t, err := time.Parse(time.RFC3339, c.Timestamp); err == nil {
				committedAt = &t
			}
			commit := &entity.IssueCommit{
				ID:          uuid.NewString(),
				IssueID:     issue.ID,
				RepoID:      ghRepo.ID,
				SHA:         c.ID,
				Message:     c.Message,
				AuthorName:  c.Author.Name,
				AuthorEmail: c.Author.Email,
				CommittedAt: committedAt,
				URL:         c.URL,
				CreatedAt:   time.Now().UTC(),
			}
			_ = uc.commitRepo.Create(ctx, commit)
		}
	}
	return nil
}

func (uc *useCase) HandlePREvent(ctx context.Context, repoFullName string, ev *ghinfra.PREvent) error {
	ghRepo, err := uc.repoRepo.GetByRepoFullName(ctx, repoFullName)
	if err != nil {
		return nil
	}

	pr := ev.PullRequest
	text := pr.Title + " " + pr.Body
	keys := ghinfra.ExtractIssueKeys(text)
	closingKeys := ghinfra.ExtractClosingKeys(text)

	for _, key := range keys {
		issue, err := uc.issueRepo.GetByKey(ctx, key)
		if err != nil {
			continue
		}

		issuePR := &entity.IssuePullRequest{
			ID:          uuid.NewString(),
			IssueID:     issue.ID,
			RepoID:      ghRepo.ID,
			PRNumber:    ev.Number,
			Title:       pr.Title,
			State:       pr.State,
			URL:         pr.HTMLURL,
			AuthorLogin: pr.User.Login,
			CreatedAt:   time.Now().UTC(),
		}
		_ = uc.prRepo.Create(ctx, issuePR)

		// PR state yangilash (agar closed/merged bo'lsa)
		if ev.Action == "closed" || ev.Action == "merged" {
			_ = uc.prRepo.UpdateState(ctx, issue.ID, ev.Number, pr.State, pr.MergedAt)
		}
	}

	// PR merged + closing keywords → issue'ni done statusiga o'tkazish
	if ev.Action == "closed" && pr.MergedAt != nil && len(closingKeys) > 0 {
		for _, key := range closingKeys {
			issue, err := uc.issueRepo.GetByKey(ctx, key)
			if err != nil {
				continue
			}
			project, err := uc.projectRepo.GetByID(ctx, issue.ProjectID)
			if err != nil {
				continue
			}
			wf, err := uc.workflowRepo.GetWithDetails(ctx, project.WorkflowID)
			if err != nil {
				continue
			}
			var doneStatusID string
			for _, s := range wf.Statuses {
				if s.Category == "done" {
					doneStatusID = s.ID
					break
				}
			}
			if doneStatusID != "" {
				_ = uc.issueRepo.UpdateStatus(ctx, issue.ID, doneStatusID)
			}
		}
	}

	return nil
}

func (uc *useCase) ListCommits(ctx context.Context, issueID string) ([]*entity.IssueCommit, error) {
	return uc.commitRepo.ListByIssue(ctx, issueID)
}

func (uc *useCase) ListPRs(ctx context.Context, issueID string) ([]*entity.IssuePullRequest, error) {
	return uc.prRepo.ListByIssue(ctx, issueID)
}
