package entity

import "time"

type Invite struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	TokenHash  string     `json:"-"`
	InvitedBy  string     `json:"invited_by"`
	ExpiresAt  time.Time  `json:"expires_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	AcceptedBy *string    `json:"accepted_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`

	Inviter *UserShort `json:"inviter,omitempty"`
}

type CreateInviteReq struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role"  validate:"required,oneof=admin member viewer"`
}

type AcceptInviteReq struct {
	Token    string `json:"token"     validate:"required"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Password string `json:"password"  validate:"required,min=8,max=72"`
}
