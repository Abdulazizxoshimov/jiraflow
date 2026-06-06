package entity

import "time"

const (
	ActivityActionCreated     = "created"
	ActivityActionUpdated     = "updated"
	ActivityActionDeleted     = "deleted"
	ActivityActionCommented   = "commented"
	ActivityActionTransitioned = "transitioned"
	ActivityActionLinked      = "linked"
	ActivityActionMentioned   = "mentioned"

	ActivityEntityIssue   = "issue"
	ActivityEntityPage    = "page"
	ActivityEntityComment = "comment"
	ActivityEntitySprint  = "sprint"
	ActivityEntitySpace   = "space"
	ActivityEntityProject = "project"
)

type ActivityEvent struct {
	ID          string         `json:"id"`
	ActorID     string         `json:"actor_id"`
	Action      string         `json:"action"`
	EntityType  string         `json:"entity_type"`
	EntityID    string         `json:"entity_id"`
	EntityTitle string         `json:"entity_title"`
	ProjectID   *string        `json:"project_id,omitempty"`
	SpaceID     *string        `json:"space_id,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`

	Actor *UserShort `json:"actor,omitempty"`
}

type ActivityFilter struct {
	Filter
	ActorID    string `form:"actor_id"    json:"actor_id,omitempty"`
	ProjectID  string `form:"project_id"  json:"project_id,omitempty"`
	SpaceID    string `form:"space_id"    json:"space_id,omitempty"`
	EntityType string `form:"entity_type" json:"entity_type,omitempty"`
}
