package websocket

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	TypeNotification MessageType = "notification"
	TypeIssueUpdated MessageType = "issue.updated"
	TypeIssueMoved   MessageType = "issue.moved"
	TypeCommentAdded MessageType = "comment.created"
	TypeMention      MessageType = "mention"
	TypeSprintUpdate MessageType = "sprint.updated"
	TypePageLocked   MessageType = "page.locked"
	TypePageUnlocked MessageType = "page.unlocked"
	TypePageUpdated  MessageType = "page.updated"
	TypeSubscribe    MessageType = "subscribe"
	TypeUnsubscribe  MessageType = "unsubscribe"
)

// Message is the envelope sent over the WebSocket connection.
type Message struct {
	Type      MessageType     `json:"type"`
	Room      string          `json:"room,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// Send delivers msg to all connections belonging to userID.
// Slow clients are skipped (non-blocking).
func (h *Hub) Send(userID string, msg Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	conns := h.clients[userID]
	h.mu.RUnlock()

	for c := range conns {
		select {
		case c.send <- b:
		default:
		}
	}
}

// Broadcast delivers msg to every connected client.
func (h *Hub) Broadcast(msg Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conns := range h.clients {
		for c := range conns {
			select {
			case c.send <- b:
			default:
			}
		}
	}
}

func newMsg(msgType MessageType, room string, payload any) Message {
	raw, _ := json.Marshal(payload)
	return Message{Type: msgType, Room: room, Payload: raw, CreatedAt: time.Now().UTC()}
}

func NewNotificationMsg(userID string, payload any) Message {
	return newMsg(TypeNotification, "user:"+userID, payload)
}

func NewIssueUpdatedMsg(projectID string, payload any) Message {
	return newMsg(TypeIssueUpdated, "project:"+projectID, payload)
}

func NewIssueMovedMsg(projectID string, payload any) Message {
	return newMsg(TypeIssueMoved, "project:"+projectID, payload)
}

func NewCommentAddedMsg(issueID string, payload any) Message {
	return newMsg(TypeCommentAdded, "issue:"+issueID, payload)
}

func NewMentionMsg(userID string, payload any) Message {
	return newMsg(TypeMention, "user:"+userID, payload)
}

func NewPageLockedMsg(pageID string, payload any) Message {
	return newMsg(TypePageLocked, "page:"+pageID, payload)
}

func NewPageUnlockedMsg(pageID string, payload any) Message {
	return newMsg(TypePageUnlocked, "page:"+pageID, payload)
}

func NewPageUpdatedMsg(spaceID string, payload any) Message {
	return newMsg(TypePageUpdated, "space:"+spaceID, payload)
}
