package entity

import "time"

// Component — loyiha ichidagi tizim qismi (frontend, backend, mobile va h.k.).
type Component struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	LeadID      *string    `json:"lead_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`

	Lead *UserShort `json:"lead,omitempty"`
}

type CreateComponentReq struct {
	Name        string  `json:"name"        validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	LeadID      *string `json:"lead_id"     validate:"omitempty,uuid4"`
}

type UpdateComponentReq struct {
	Name        *string `json:"name"        validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	LeadID      *string `json:"lead_id"     validate:"omitempty,uuid4"`
}
