package entity

import "time"

type ProjectTemplate struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	Type                  string         `json:"type"` // scrum | kanban | business
	Description           *string        `json:"description,omitempty"`
	IconURL               *string        `json:"icon_url,omitempty"`
	DefaultWorkflowConfig map[string]any `json:"default_workflow_config,omitempty"`
	IsSystem              bool           `json:"is_system"`
	CreatedAt             time.Time      `json:"created_at"`
}
