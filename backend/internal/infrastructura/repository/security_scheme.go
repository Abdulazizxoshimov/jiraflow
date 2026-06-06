package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SecuritySchemeRepository interface {
	Create(ctx context.Context, req *entity.CreateSecuritySchemeReq) (*entity.SecurityScheme, error)
	GetByID(ctx context.Context, id string) (*entity.SecurityScheme, error)
	List(ctx context.Context, projectID string) ([]*entity.SecurityScheme, error)
	Delete(ctx context.Context, id string) error

	AddLevel(ctx context.Context, schemeID string, req *entity.CreateSecurityLevelReq) (*entity.SecurityLevel, error)
	GetLevel(ctx context.Context, levelID string) (*entity.SecurityLevel, error)
	DeleteLevel(ctx context.Context, levelID string) error

	AddMember(ctx context.Context, levelID string, req *entity.CreateSecurityLevelMemberReq) (*entity.SecurityLevelMember, error)
	DeleteMember(ctx context.Context, memberID string) error
}
