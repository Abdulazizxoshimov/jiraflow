package entity

import "time"

type PageVersion struct {
	ID          string         `json:"id"`
	PageID      string         `json:"page_id"`
	Version     int            `json:"version"`
	Title       string         `json:"title"`
	Content     map[string]any `json:"content"`
	ContentText string         `json:"content_text"`
	AuthorID    string         `json:"author_id"`
	ChangeNote  *string        `json:"change_note,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`

	Author *UserShort `json:"author,omitempty"`
}

// DiffLine represents one line in a text diff.
type DiffLine struct {
	Op   string `json:"op"`   // "equal" | "insert" | "delete"
	Text string `json:"text"`
}

// PageVersionDiff is the result of comparing two page versions.
type PageVersionDiff struct {
	PageID   string      `json:"page_id"`
	V1       int         `json:"v1"`
	V2       int         `json:"v2"`
	TitleV1  string      `json:"title_v1"`
	TitleV2  string      `json:"title_v2"`
	Lines    []DiffLine  `json:"lines"`
}
