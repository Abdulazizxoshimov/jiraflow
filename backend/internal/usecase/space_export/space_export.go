package space_export

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/minio"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/tiptap"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo     repository.SpaceExportRepository
	pageRepo repository.PageRepository
	minio    minio.Client
	log      logger.Logger
}

func New(
	repo repository.SpaceExportRepository,
	pageRepo repository.PageRepository,
	minio minio.Client,
	log logger.Logger,
) UseCase {
	return &useCase{repo: repo, pageRepo: pageRepo, minio: minio, log: log}
}

func (uc *useCase) RequestExport(ctx context.Context, spaceID, requestedBy string) (*entity.SpaceExport, error) {
	e := &entity.SpaceExport{
		ID:          uuid.NewString(),
		SpaceID:     spaceID,
		RequestedBy: requestedBy,
		Status:      "pending",
	}
	if err := uc.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("space_export.RequestExport: %w", err)
	}
	go uc.processExport(context.Background(), e.ID, spaceID)
	return e, nil
}

func (uc *useCase) GetExport(ctx context.Context, id string) (*entity.SpaceExport, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) ListExports(ctx context.Context, spaceID string) ([]*entity.SpaceExport, error) {
	return uc.repo.ListBySpace(ctx, spaceID)
}

func (uc *useCase) processExport(ctx context.Context, exportID, spaceID string) {
	processing := "processing"
	_ = uc.repo.UpdateStatus(ctx, exportID, processing, nil, nil)

	fileURL, err := uc.buildAndUploadZip(ctx, spaceID)
	if err != nil {
		uc.log.Error(ctx, "space_export.processExport failed",
			logger.String("export_id", exportID), logger.SafeString("err", err.Error()))
		errMsg := err.Error()
		_ = uc.repo.UpdateStatus(ctx, exportID, "failed", nil, &errMsg)
		return
	}

	_ = uc.repo.UpdateStatus(ctx, exportID, "done", &fileURL, nil)
	uc.log.Info(ctx, "space export done", logger.String("export_id", exportID))
}

func (uc *useCase) buildAndUploadZip(ctx context.Context, spaceID string) (string, error) {
	pages, _, err := uc.pageRepo.List(ctx, &entity.PageFilter{SpaceID: spaceID, Status: "published"})
	if err != nil {
		return "", fmt.Errorf("buildAndUploadZip list pages: %w", err)
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, p := range pages {
		var bodyHTML string
		if len(p.Content) > 0 {
			contentJSON, _ := json.Marshal(p.Content)
			doc, parseErr := tiptap.Parse(contentJSON)
			if parseErr == nil {
				bodyHTML = doc.RenderHTML()
			} else {
				bodyHTML = "<p>" + template.HTMLEscapeString(p.ContentText) + "</p>"
			}
		}

		htmlContent := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>%s</title></head><body><h1>%s</h1>%s</body></html>`,
			template.HTMLEscapeString(p.Title), template.HTMLEscapeString(p.Title), bodyHTML)

		filename := path.Join("pages", sanitizeExportFilename(p.Title)+".html")
		f, err := zw.Create(filename)
		if err != nil {
			continue
		}
		_, _ = f.Write([]byte(htmlContent))
	}

	if err := zw.Close(); err != nil {
		return "", fmt.Errorf("buildAndUploadZip zip close: %w", err)
	}

	objectName := fmt.Sprintf("space-exports/%s/%s.zip", spaceID, time.Now().UTC().Format("20060102-150405"))
	data := buf.Bytes()
	if _, err := uc.minio.Upload(ctx, objectName, "application/zip", bytes.NewReader(data), int64(len(data))); err != nil {
		return "", fmt.Errorf("buildAndUploadZip upload: %w", err)
	}

	url, err := uc.minio.PresignedURL(ctx, objectName, 72*time.Hour)
	if err != nil {
		return "", fmt.Errorf("buildAndUploadZip presign: %w", err)
	}
	return url, nil
}

func sanitizeExportFilename(title string) string {
	var b bytes.Buffer
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		} else if r == ' ' {
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "page"
	}
	return b.String()
}
