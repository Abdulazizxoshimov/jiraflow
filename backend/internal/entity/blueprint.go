package entity

import "time"

type Blueprint struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  *string        `json:"description,omitempty"`
	IconURL      *string        `json:"icon_url,omitempty"`
	Category     *string        `json:"category,omitempty"`
	TemplateBody *string        `json:"template_body,omitempty"`
	Schema       map[string]any `json:"schema,omitempty"`
	IsSystem     bool           `json:"is_system"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreateBlueprintReq struct {
	Name         string         `json:"name"         validate:"required,min=1,max=255"`
	Description  *string        `json:"description"`
	IconURL      *string        `json:"icon_url"`
	Category     *string        `json:"category"     validate:"omitempty,max=100"`
	TemplateBody *string        `json:"template_body"`
	Schema       map[string]any `json:"schema"`
}

type CreatePageFromBlueprintReq struct {
	SpaceID  string  `json:"space_id"  validate:"required,uuid4"`
	ParentID *string `json:"parent_id" validate:"omitempty,uuid4"`
	Title    string  `json:"title"     validate:"required,min=1,max=500"`
}
