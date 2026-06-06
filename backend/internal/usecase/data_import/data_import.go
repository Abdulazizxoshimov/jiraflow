package data_import

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type useCase struct {
	importRepo   repository.DataImportRepository
	issueRepo    repository.IssueRepository
	projectRepo  repository.ProjectRepository
	workflowRepo repository.WorkflowRepository
}

func New(
	importRepo repository.DataImportRepository,
	issueRepo repository.IssueRepository,
	projectRepo repository.ProjectRepository,
	workflowRepo repository.WorkflowRepository,
) UseCase {
	return &useCase{
		importRepo:   importRepo,
		issueRepo:    issueRepo,
		projectRepo:  projectRepo,
		workflowRepo: workflowRepo,
	}
}

func (uc *useCase) GetStatus(ctx context.Context, id string) (*entity.DataImport, error) {
	return uc.importRepo.GetByID(ctx, id)
}

func (uc *useCase) ImportJira(ctx context.Context, userID string, data []byte) (*entity.DataImport, error) {
	imp := uc.createJob(ctx, userID, "jira")
	go uc.runJiraImport(imp.ID, userID, data)
	return imp, nil
}

func (uc *useCase) ImportTrello(ctx context.Context, userID string, data []byte) (*entity.DataImport, error) {
	imp := uc.createJob(ctx, userID, "trello")
	go uc.runTrelloImport(imp.ID, userID, data)
	return imp, nil
}

func (uc *useCase) ImportLinear(ctx context.Context, userID string, data []byte) (*entity.DataImport, error) {
	imp := uc.createJob(ctx, userID, "linear")
	go uc.runLinearImport(imp.ID, userID, data)
	return imp, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (uc *useCase) createJob(ctx context.Context, userID, source string) *entity.DataImport {
	imp := &entity.DataImport{
		ID:        uuid.NewString(),
		UserID:    userID,
		Source:    source,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
	}
	_ = uc.importRepo.Create(ctx, imp)
	return imp
}

func (uc *useCase) defaultStatusID(ctx context.Context) string {
	wf, err := uc.workflowRepo.GetDefault(ctx)
	if err != nil || wf == nil {
		return ""
	}
	statuses, err := uc.workflowRepo.ListStatuses(ctx, wf.ID)
	if err != nil || len(statuses) == 0 {
		return ""
	}
	return statuses[0].ID
}

func (uc *useCase) defaultProjectID(ctx context.Context) string {
	projects, _, err := uc.projectRepo.List(ctx, &entity.ProjectFilter{Filter: entity.Filter{Page: 1, Limit: 1}})
	if err != nil || len(projects) == 0 {
		return ""
	}
	return projects[0].ID
}

func (uc *useCase) runJiraImport(importID, userID string, data []byte) {
	ctx := context.Background()
	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", 0, 0, "")

	var export entity.JiraXMLExport
	if err := xml.Unmarshal(data, &export); err != nil {
		_ = uc.importRepo.UpdateStatus(ctx, importID, "failed", 0, 0, "invalid jira xml: "+err.Error())
		return
	}

	statusID := uc.defaultStatusID(ctx)
	projectID := uc.defaultProjectID(ctx)
	total := len(export.Issues)
	processed := 0

	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", total, processed, "")

	for _, ji := range export.Issues {
		desc := ji.Description
		issue := &entity.Issue{
			ID:          uuid.NewString(),
			ProjectID:   projectID,
			Title:       ji.Title,
			Description: &desc,
			Type:        normalizeIssueType(ji.Type),
			Priority:    normalizePriority(ji.Priority),
			StatusID:    statusID,
			ReporterID:  userID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		_ = uc.issueRepo.Create(ctx, issue)
		processed++
		if processed%10 == 0 {
			_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", total, processed, "")
		}
	}

	_ = uc.importRepo.UpdateStatus(ctx, importID, "done", total, processed, "")
	_ = uc.importRepo.MarkCompleted(ctx, importID)
}

func (uc *useCase) runTrelloImport(importID, userID string, data []byte) {
	ctx := context.Background()
	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", 0, 0, "")

	var export entity.TrelloExport
	if err := json.Unmarshal(data, &export); err != nil {
		_ = uc.importRepo.UpdateStatus(ctx, importID, "failed", 0, 0, "invalid trello json: "+err.Error())
		return
	}

	// Build list name index
	listNames := make(map[string]string)
	for _, l := range export.Lists {
		listNames[l.ID] = l.Name
	}

	statusID := uc.defaultStatusID(ctx)
	projectID := uc.defaultProjectID(ctx)

	cards := make([]entity.TrelloCard, 0, len(export.Cards))
	for _, c := range export.Cards {
		if !c.Closed {
			cards = append(cards, c)
		}
	}

	total := len(cards)
	processed := 0
	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", total, processed, "")

	for _, card := range cards {
		desc := card.Desc
		issue := &entity.Issue{
			ID:          uuid.NewString(),
			ProjectID:   projectID,
			Title:       card.Name,
			Description: &desc,
			Type:        "task",
			Priority:    "medium",
			StatusID:    statusID,
			ReporterID:  userID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		_ = uc.issueRepo.Create(ctx, issue)
		processed++
	}

	_ = uc.importRepo.UpdateStatus(ctx, importID, "done", total, processed, "")
	_ = uc.importRepo.MarkCompleted(ctx, importID)
}

func (uc *useCase) runLinearImport(importID, userID string, data []byte) {
	ctx := context.Background()
	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", 0, 0, "")

	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		_ = uc.importRepo.UpdateStatus(ctx, importID, "failed", 0, 0, "invalid linear csv: "+err.Error())
		return
	}

	if len(records) < 2 {
		_ = uc.importRepo.UpdateStatus(ctx, importID, "done", 0, 0, "")
		_ = uc.importRepo.MarkCompleted(ctx, importID)
		return
	}

	// Build header index
	headers := records[0]
	idx := make(map[string]int)
	for i, h := range headers {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	statusID := uc.defaultStatusID(ctx)
	projectID := uc.defaultProjectID(ctx)

	rows := records[1:]
	total := len(rows)
	processed := 0
	_ = uc.importRepo.UpdateStatus(ctx, importID, "processing", total, processed, "")

	get := func(row []string, col string) string {
		if i, ok := idx[col]; ok && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}

	for _, row := range rows {
		title := get(row, "title")
		if title == "" {
			title = get(row, "name")
		}
		if title == "" {
			processed++
			continue
		}
		desc := fmt.Sprintf("%s\n\nImported from Linear.", get(row, "description"))
		issue := &entity.Issue{
			ID:          uuid.NewString(),
			ProjectID:   projectID,
			Title:       title,
			Description: &desc,
			Type:        "task",
			Priority:    normalizePriority(get(row, "priority")),
			StatusID:    statusID,
			ReporterID:  userID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		_ = uc.issueRepo.Create(ctx, issue)
		processed++
	}

	_ = uc.importRepo.UpdateStatus(ctx, importID, "done", total, processed, "")
	_ = uc.importRepo.MarkCompleted(ctx, importID)
}

func normalizeIssueType(t string) string {
	switch strings.ToLower(t) {
	case "bug":
		return "bug"
	case "story", "user story":
		return "story"
	case "epic":
		return "epic"
	case "sub-task", "subtask":
		return "subtask"
	default:
		return "task"
	}
}

func normalizePriority(p string) string {
	switch strings.ToLower(p) {
	case "highest", "critical", "urgent":
		return "highest"
	case "high":
		return "high"
	case "low":
		return "low"
	case "lowest", "trivial":
		return "lowest"
	default:
		return "medium"
	}
}
