package entity

import "time"

type IssueType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IconURL     *string   `json:"icon_url,omitempty"`
	Color       *string   `json:"color,omitempty"`
	IsSubtask   bool      `json:"is_subtask"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
}

type IssueTypeScheme struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	ProjectID  *string      `json:"project_id,omitempty"`
	IssueTypes []*IssueType `json:"issue_types,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}

type CreateIssueTypeReq struct {
	Name        string  `json:"name"        validate:"required,min=1,max=100"`
	Description *string `json:"description"`
	IconURL     *string `json:"icon_url"`
	Color       *string `json:"color"       validate:"omitempty,max=20"`
	IsSubtask   bool    `json:"is_subtask"`
}

type CreateIssueTypeSchemeReq struct {
	Name         string   `json:"name"          validate:"required,min=1,max=255"`
	ProjectID    *string  `json:"project_id"    validate:"omitempty,uuid4"`
	IssueTypeIDs []string `json:"issue_type_ids" validate:"required,min=1"`
}
