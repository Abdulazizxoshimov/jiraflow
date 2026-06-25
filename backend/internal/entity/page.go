package entity

import "time"

type Page struct {
	ID             string         `json:"id"`
	SpaceID        string         `json:"space_id"`
	ParentID       *string        `json:"parent_id,omitempty"`
	Title          string         `json:"title"`
	Icon           string         `json:"icon"`
	Content        map[string]any `json:"content"`
	ContentText    string         `json:"content_text"`
	AuthorID       string         `json:"author_id"`
	LastEditorID   string         `json:"last_editor_id"`
	CurrentVersion int            `json:"current_version"`
	Status         string         `json:"status"` // draft | published
	Position       int            `json:"position"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      *time.Time     `json:"-"`

	Author     *UserShort `json:"author,omitempty"`
	LastEditor *UserShort `json:"last_editor,omitempty"`
	Children   []Page     `json:"children,omitempty"`
}

// PageTree is a lightweight node used when building the space page tree.
type PageTree struct {
	ID       string     `json:"id"`
	ParentID *string    `json:"parent_id,omitempty"`
	Title    string     `json:"title"`
	Icon     string     `json:"icon"`
	Position int        `json:"position"`
	Status   string     `json:"status"`
	Children []PageTree `json:"children,omitempty"`
}

type CreatePageReq struct {
	ParentID    *string        `json:"parent_id"    validate:"omitempty,uuid4"`
	Title       string         `json:"title"        validate:"required,min=1,max=500"`
	Content     map[string]any `json:"content"`
	ContentText string         `json:"content_text"`
	Status      string         `json:"status"       validate:"omitempty,oneof=draft published"`
	ChangeNote  *string        `json:"change_note"  validate:"omitempty,max=500"`
}

type UpdatePageReq struct {
	Title       *string        `json:"title"        validate:"omitempty,min=1,max=500"`
	Icon        *string        `json:"icon"         validate:"omitempty,max=10"`
	Content     map[string]any `json:"content"`
	ContentText *string        `json:"content_text"`
	Status      *string        `json:"status"       validate:"omitempty,oneof=draft published"`
	ChangeNote  *string        `json:"change_note"  validate:"omitempty,max=500"`
}

type PageFilter struct {
	Filter
	SpaceID       string `form:"space_id"  json:"space_id"`
	ParentID      string `form:"parent_id" json:"parent_id"`
	Status        string `form:"status"    json:"status" validate:"omitempty,oneof=draft published"`
	AuthorID      string `form:"author_id" json:"author_id"`
	CQL           string `form:"cql"       json:"cql"`
	CurrentUserID string `form:"-"         json:"-"`
}

type PageWatcher struct {
	PageID    string    `json:"page_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User *UserShort `json:"user,omitempty"`
}

type CopyPageReq struct {
	Title        string  `json:"title"        validate:"required,min=1,max=500"`
	NewSpaceID   *string `json:"space_id"     validate:"omitempty,uuid4"`
	NewParentID  *string `json:"parent_id"    validate:"omitempty,uuid4"`
	CopyChildren bool    `json:"copy_children"`
}
