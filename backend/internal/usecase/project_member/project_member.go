package project_member

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	memberRepo  repository.ProjectMemberRepository
	projectRepo repository.ProjectRepository
	log         logger.Logger
}

func New(memberRepo repository.ProjectMemberRepository, projectRepo repository.ProjectRepository, log logger.Logger) UseCase {
	return &useCase{memberRepo: memberRepo, projectRepo: projectRepo, log: log}
}

func (uc *useCase) Add(ctx context.Context, projectID string, req *entity.AddProjectMemberReq, actorID string) error {
	if _, err := uc.projectRepo.GetByID(ctx, projectID); err != nil {
		return apperr.NotFound("project")
	}
	isMember, err := uc.memberRepo.IsMember(ctx, projectID, req.UserID)
	if err != nil {
		return err
	}
	if isMember {
		return apperr.Conflict("user is already a project member")
	}
	m := &entity.ProjectMember{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.memberRepo.Add(ctx, m); err != nil {
		uc.log.Error(ctx, "project_member.Add: db error", logger.String("project_id", projectID), logger.String("user_id", req.UserID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "project member added", logger.String("project_id", projectID), logger.String("user_id", req.UserID))
	return nil
}

func (uc *useCase) UpdateRole(ctx context.Context, projectID, userID string, req *entity.UpdateProjectMemberRoleReq, actorID string) error {
	actor, err := uc.memberRepo.GetMember(ctx, projectID, actorID)
	if err != nil || actor.Role != "admin" {
		uc.log.Warn(ctx, "project_member.UpdateRole: forbidden", logger.String("project_id", projectID), logger.String("actor_id", actorID))
		return apperr.Forbidden("only project admin can change member roles")
	}
	if err := uc.memberRepo.UpdateRole(ctx, projectID, userID, req.Role); err != nil {
		uc.log.Error(ctx, "project_member.UpdateRole: db error", logger.String("project_id", projectID), logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "project member role updated", logger.String("project_id", projectID), logger.String("user_id", userID), logger.String("role", req.Role))
	return nil
}

func (uc *useCase) Remove(ctx context.Context, projectID, userID, actorID string) error {
	actor, err := uc.memberRepo.GetMember(ctx, projectID, actorID)
	if err != nil || actor.Role != "admin" {
		uc.log.Warn(ctx, "project_member.Remove: forbidden", logger.String("project_id", projectID), logger.String("actor_id", actorID))
		return apperr.Forbidden("only project admin can remove members")
	}
	if err := uc.memberRepo.Remove(ctx, projectID, userID); err != nil {
		uc.log.Error(ctx, "project_member.Remove: db error", logger.String("project_id", projectID), logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "project member removed", logger.String("project_id", projectID), logger.String("user_id", userID))
	return nil
}

func (uc *useCase) ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.ProjectMember, int, error) {
	return uc.memberRepo.ListByProject(ctx, projectID, filter)
}

func (uc *useCase) IsMember(ctx context.Context, projectID, userID string) (bool, error) {
	return uc.memberRepo.IsMember(ctx, projectID, userID)
}
