package page_version

import (
	"context"
	"fmt"
	"strings"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	versionRepo repository.PageVersionRepository
	log         logger.Logger
}

func New(versionRepo repository.PageVersionRepository, log logger.Logger) UseCase {
	return &useCase{versionRepo: versionRepo, log: log}
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.PageVersion, error) {
	return uc.versionRepo.GetByID(ctx, id)
}

func (uc *useCase) GetByVersion(ctx context.Context, pageID string, version int) (*entity.PageVersion, error) {
	return uc.versionRepo.GetByVersion(ctx, pageID, version)
}

func (uc *useCase) ListByPage(ctx context.Context, pageID string, filter *entity.Filter) ([]*entity.PageVersion, int, error) {
	return uc.versionRepo.ListByPage(ctx, pageID, filter)
}

// Diff computes a line-level diff between two page versions using Myers diff.
func (uc *useCase) Diff(ctx context.Context, pageID string, v1, v2 int) (*entity.PageVersionDiff, error) {
	pv1, err := uc.versionRepo.GetByVersion(ctx, pageID, v1)
	if err != nil {
		return nil, fmt.Errorf("page_version.Diff v%d: %w", v1, err)
	}
	pv2, err := uc.versionRepo.GetByVersion(ctx, pageID, v2)
	if err != nil {
		return nil, fmt.Errorf("page_version.Diff v%d: %w", v2, err)
	}

	lines := diffLines(pv1.ContentText, pv2.ContentText)
	return &entity.PageVersionDiff{
		PageID:  pageID,
		V1:      v1,
		V2:      v2,
		TitleV1: pv1.Title,
		TitleV2: pv2.Title,
		Lines:   lines,
	}, nil
}

// diffLines performs a simple line-level LCS diff between two texts.
func diffLines(a, b string) []entity.DiffLine {
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")

	// Build LCS table
	m, n := len(linesA), len(linesB)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if linesA[i] == linesB[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	var result []entity.DiffLine
	i, j := 0, 0
	for i < m || j < n {
		switch {
		case i < m && j < n && linesA[i] == linesB[j]:
			result = append(result, entity.DiffLine{Op: "equal", Text: linesA[i]})
			i++
			j++
		case j < n && (i >= m || dp[i][j+1] >= dp[i+1][j]):
			result = append(result, entity.DiffLine{Op: "insert", Text: linesB[j]})
			j++
		default:
			result = append(result, entity.DiffLine{Op: "delete", Text: linesA[i]})
			i++
		}
	}
	return result
}
