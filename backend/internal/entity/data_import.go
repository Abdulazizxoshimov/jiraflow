package entity

import "time"

type DataImport struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Source         string     `json:"source"` // jira | trello | linear
	Status         string     `json:"status"` // pending | processing | done | failed
	TotalItems     int        `json:"total_items"`
	ProcessedItems int        `json:"processed_items"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

// JiraXMLExport top-level structure for Jira XML backup.
type JiraXMLExport struct {
	Issues []JiraXMLIssue `xml:"channel>item"`
}

type JiraXMLIssue struct {
	Title       string `xml:"title"`
	Key         string `xml:"key"`
	Type        string `xml:"type"`
	Priority    string `xml:"priority"`
	Status      string `xml:"status"`
	Description string `xml:"description"`
	Assignee    string `xml:"assignee"`
	Reporter    string `xml:"reporter"`
}

// TrelloExport is the structure of Trello's JSON export.
type TrelloExport struct {
	Name  string        `json:"name"`
	Lists []TrelloList  `json:"lists"`
	Cards []TrelloCard  `json:"cards"`
}

type TrelloList struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Closed bool   `json:"closed"`
}

type TrelloCard struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Desc        string   `json:"desc"`
	IDList      string   `json:"idList"`
	Closed      bool     `json:"closed"`
	Labels      []struct{ Name string `json:"name"` } `json:"labels"`
}

// LinearExport represents a row in Linear's CSV issue export.
type LinearExport struct {
	ID          string
	Title       string
	Description string
	Status      string
	Priority    string
	Assignee    string
	Labels      string
	CreatedAt   string
}
