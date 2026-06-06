package entity

import "time"

type PageTemplate struct {
	ID          string         `json:"id"`
	SpaceID     *string        `json:"space_id,omitempty"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Category    string         `json:"category"` // general | meeting | retrospective | decision
	Content     map[string]any `json:"content"`
	ContentText string         `json:"content_text"`
	Icon        *string        `json:"icon,omitempty"`
	CreatedBy   string         `json:"created_by"`
	IsGlobal    bool           `json:"is_global"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreatePageTemplateReq struct {
	Name        string         `json:"name"        binding:"required,min=1,max=255"`
	Description *string        `json:"description"`
	Category    string         `json:"category"    binding:"required,oneof=general meeting retrospective decision"`
	Content     map[string]any `json:"content"`
	ContentText string         `json:"content_text"`
	Icon        *string        `json:"icon"`
	IsGlobal    bool           `json:"is_global"`
}

type UpdatePageTemplateReq struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Category    *string        `json:"category"`
	Content     map[string]any `json:"content"`
	ContentText *string        `json:"content_text"`
	Icon        *string        `json:"icon"`
}

type PageTemplateFilter struct {
	SpaceID  string
	Category string
	IsGlobal *bool
	Filter
}
