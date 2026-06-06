package entity

import "time"

type SearchResultType string

const (
	SearchResultIssue   SearchResultType = "issue"
	SearchResultPage    SearchResultType = "page"
	SearchResultProject SearchResultType = "project"
	SearchResultSpace   SearchResultType = "space"
)

type SearchResult struct {
	Type      SearchResultType `json:"type"`
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Excerpt   string           `json:"excerpt,omitempty"`
	URL       string           `json:"url,omitempty"`
	UpdatedAt time.Time        `json:"updated_at"`
	Meta      map[string]any   `json:"meta,omitempty"`
}

type SearchFilter struct {
	Query     string             `form:"q"    json:"q"    validate:"required,min=1"`
	Types     []SearchResultType `form:"type" json:"type"`
	ProjectID string             `form:"project_id" json:"project_id,omitempty"`
	SpaceID   string             `form:"space_id"   json:"space_id,omitempty"`
	Page      int                `form:"page"  json:"page"`
	Limit     int                `form:"limit" json:"limit"`
}

type SearchSuggestion struct {
	ID    string           `json:"id"`
	Type  SearchResultType `json:"type"`
	Title string           `json:"title"`
}
