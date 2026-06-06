package entity

import "time"

type PermissionScheme struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	CreatedBy   string                   `json:"created_by"`
	Grants      []*PermissionSchemeGrant `json:"grants,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type PermissionSchemeGrant struct {
	ID         string    `json:"id"`
	SchemeID   string    `json:"scheme_id"`
	Permission string    `json:"permission"`
	HolderType string    `json:"holder_type"` // user | role | anyone
	HolderID   *string   `json:"holder_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreatePermissionSchemeReq struct {
	Name        string `json:"name"        validate:"required,max=100"`
	Description string `json:"description"`
}

type UpdatePermissionSchemeReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type AddGrantReq struct {
	Permission string  `json:"permission"   validate:"required"`
	HolderType string  `json:"holder_type"  validate:"required,oneof=user role anyone"`
	HolderID   *string `json:"holder_id"`
}

type AssignSchemeReq struct {
	SchemeID string `json:"scheme_id" validate:"required"`
}
