package entity

import "time"

type PageMacro struct {
	ID        string         `json:"id"`
	PageID    string         `json:"page_id"`
	MacroType string         `json:"macro_type"`
	Config    map[string]any `json:"config"`
	CreatedAt time.Time      `json:"created_at"`
}

type UpsertPageMacroReq struct {
	MacroType string         `json:"macro_type" binding:"required"`
	Config    map[string]any `json:"config"     binding:"required"`
}

const (
	MacroTypeJiraIssue = "jira_issue"
	MacroTypeChart     = "chart"
	MacroTypeStatus    = "status"
)
