package entity

import "time"

type ProjectMember struct {
	ProjectID string    `json:"project_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type AddProjectMemberReq struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Role   string `json:"role"    validate:"required,oneof=admin member viewer"`
}

type UpdateProjectMemberRoleReq struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}
