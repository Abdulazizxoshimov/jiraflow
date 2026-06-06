package invite

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateInviteReq, invitedBy string) (*entity.Invite, error)
	Accept(ctx context.Context, req *entity.AcceptInviteReq) (*entity.TokenPair, error)
	ListPending(ctx context.Context, filter *entity.Filter) ([]*entity.Invite, int, error)
	Revoke(ctx context.Context, id string) error
}
