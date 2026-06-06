package entity

import "time"

type Issue struct {
	ID                 string         `json:"id"`
	ProjectID          string         `json:"project_id"`
	IssueNumber        int            `json:"issue_number"`
	Title              string         `json:"title"`
	Description        *string        `json:"description,omitempty"`
	Type               string         `json:"type"`     // task | bug | story | epic | subtask
	StatusID           string         `json:"status_id"`
	Priority           string         `json:"priority"` // lowest | low | medium | high | highest
	AssigneeID         *string        `json:"assignee_id,omitempty"`
	ReporterID         string         `json:"reporter_id"`
	ParentID           *string        `json:"parent_id,omitempty"`
	SprintID           *string        `json:"sprint_id,omitempty"`
	StoryPoints        *int           `json:"story_points,omitempty"`
	DueDate            *time.Time     `json:"due_date,omitempty"`
	OriginalEstimate   *int           `json:"original_estimate,omitempty"`  // sekundda
	RemainingEstimate  *int           `json:"remaining_estimate,omitempty"` // sekundda
	CustomFields       map[string]any `json:"custom_fields"`
	Resolution         *string        `json:"resolution,omitempty"` // fixed | wont_fix | duplicate | incomplete | cannot_reproduce | done
	Position           int            `json:"position"`
	VoteCount          int            `json:"vote_count"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          *time.Time     `json:"-"`

	Status     *WorkflowStatus `json:"status,omitempty"`
	Assignee   *UserShort      `json:"assignee,omitempty"`
	Reporter   *UserShort      `json:"reporter,omitempty"`
	Labels     []Label         `json:"labels,omitempty"`
	Components []Component     `json:"components,omitempty"`
	Versions   []Version       `json:"versions,omitempty"`        // fix versions
	AffectsVersions []Version  `json:"affects_versions,omitempty"` // affected versions
	EpicProgress *EpicProgress `json:"epic_progress,omitempty"`
}

// EpicProgress — epic'ning bajarilish foizi.
type EpicProgress struct {
	Total    int     `json:"total"`
	Done     int     `json:"done"`
	Progress float64 `json:"progress"` // 0.0–100.0
}

// IssueKey returns the human-readable key like "PROJ-42".
func (i *Issue) IssueKey(projectKey string) string {
	return projectKey + "-" + itoa(i.IssueNumber)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

type CreateIssueReq struct {
	ProjectID          string         `json:"project_id"           validate:"required,uuid4"`
	Title              string         `json:"title"                validate:"required,min=1,max=500"`
	Description        *string        `json:"description"          validate:"omitempty"`
	Type               string         `json:"type"                 validate:"required,oneof=task bug story epic subtask"`
	Priority           string         `json:"priority"             validate:"omitempty,oneof=lowest low medium high highest"`
	AssigneeID         *string        `json:"assignee_id"          validate:"omitempty,uuid4"`
	ParentID           *string        `json:"parent_id"            validate:"omitempty,uuid4"`
	SprintID           *string        `json:"sprint_id"            validate:"omitempty,uuid4"`
	StoryPoints        *int           `json:"story_points"         validate:"omitempty,gte=0"`
	DueDate            *time.Time     `json:"due_date"`
	OriginalEstimate   *int           `json:"original_estimate"    validate:"omitempty,gt=0"`
	RemainingEstimate  *int           `json:"remaining_estimate"   validate:"omitempty,gte=0"`
	LabelIDs           []string       `json:"label_ids"`
	ComponentIDs       []string       `json:"component_ids"`
	FixVersionIDs      []string       `json:"fix_version_ids"`
	AffectsVersionIDs  []string       `json:"affects_version_ids"`
	CustomFields       map[string]any `json:"custom_fields"`
}

type UpdateIssueReq struct {
	Title              *string        `json:"title"                validate:"omitempty,min=1,max=500"`
	Description        *string        `json:"description"`
	Priority           *string        `json:"priority"             validate:"omitempty,oneof=lowest low medium high highest"`
	AssigneeID         *string        `json:"assignee_id"          validate:"omitempty,uuid4"`
	SprintID           *string        `json:"sprint_id"            validate:"omitempty,uuid4"`
	StoryPoints        *int           `json:"story_points"         validate:"omitempty,gte=0"`
	DueDate            *time.Time     `json:"due_date"`
	OriginalEstimate   *int           `json:"original_estimate"    validate:"omitempty,gt=0"`
	RemainingEstimate  *int           `json:"remaining_estimate"   validate:"omitempty,gte=0"`
	Resolution         *string        `json:"resolution"           validate:"omitempty,oneof=fixed wont_fix duplicate incomplete cannot_reproduce done"`
	LabelIDs           []string       `json:"label_ids"`
	ComponentIDs       []string       `json:"component_ids"`
	FixVersionIDs      []string       `json:"fix_version_ids"`
	AffectsVersionIDs  []string       `json:"affects_version_ids"`
	CustomFields       map[string]any `json:"custom_fields"`
}

type IssueFilter struct {
	Filter
	ProjectID    string     `form:"project_id"    json:"project_id"`
	SprintID     string     `form:"sprint_id"     json:"sprint_id"`
	AssigneeID   string     `form:"assignee_id"   json:"assignee_id"`
	AssigneeIDs  []string   `form:"assignee_ids"  json:"assignee_ids"`
	ReporterID   string     `form:"reporter_id"   json:"reporter_id"`
	StatusID     string     `form:"status_id"     json:"status_id"`
	StatusIDs    []string   `form:"status_ids"    json:"status_ids"`
	Type         string     `form:"type"          json:"type"`
	Types        []string   `form:"types"         json:"types"`
	Priority     string     `form:"priority"      json:"priority"`
	Priorities   []string   `form:"priorities"    json:"priorities"`
	LabelIDs     []string   `form:"label_ids"     json:"label_ids"`
	ComponentIDs []string   `form:"component_ids" json:"component_ids"`
	VersionIDs        []string   `form:"version_ids"         json:"version_ids"`
	AffectsVersionIDs []string   `form:"affects_version_ids" json:"affects_version_ids"`
	ParentID     string     `form:"parent_id"     json:"parent_id"`
	EpicID       string     `form:"epic_id"       json:"epic_id"`
	NoSprint     bool       `form:"no_sprint"     json:"no_sprint"` // backlog
	DueDateFrom  *time.Time `form:"due_date_from" json:"due_date_from"`
	DueDateTo    *time.Time `form:"due_date_to"   json:"due_date_to"`
	CreatedAfter  *time.Time `form:"created_after"  json:"created_after"`
	CreatedBefore *time.Time `form:"created_before" json:"created_before"`
	TextSearch    string     `form:"q"              json:"q"`
	JQL           string     `form:"jql"            json:"jql"`            // Jira Query Language
	CurrentUserID string     `form:"-"              json:"-"`              // injected by handler for currentUser()
}

type TransitionIssueReq struct {
	StatusID string `json:"status_id" validate:"required,uuid4"`
}

// IssuePositionItem bitta issue va uning yangi pozitsiyasi.
type IssuePositionItem struct {
	IssueID  string `json:"issue_id"  validate:"required,uuid4"`
	Position int    `json:"position"  validate:"gte=0"`
}

// ReorderIssuesReq backlog yoki board'da drag-and-drop tartibini saqlash uchun.
type ReorderIssuesReq struct {
	Items []IssuePositionItem `json:"items" validate:"required,min=1,dive"`
}

// MoveIssueReq board'da issue'ni boshqa kolonnaga ko'chirish uchun.
type MoveIssueReq struct {
	StatusID string `json:"status_id" validate:"required,uuid4"`
	Position int    `json:"position"  validate:"gte=0"`
}
