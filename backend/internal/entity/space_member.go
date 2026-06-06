package entity

import "time"

type SpaceMember struct {
	SpaceID   string    `json:"space_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // admin | member | viewer
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type AddSpaceMemberReq struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Role   string `json:"role"    validate:"required,oneof=admin member viewer"`
}

type UpdateSpaceMemberRoleReq struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}
