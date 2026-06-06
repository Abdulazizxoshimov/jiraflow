package websocket

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	TypeNotification MessageType = "notification"
	TypeIssueUpdated MessageType = "issue_updated"
	TypeCommentAdded MessageType = "comment_added"
	TypeMention      MessageType = "mention"
	TypeSprintUpdate MessageType = "sprint_update"
)

// Message is the envelope sent over the WebSocket connection.
type Message struct {
	Type      MessageType `json:"type"`
	Payload   any         `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
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

func NewNotificationMsg(payload any) Message {
	return Message{Type: TypeNotification, Payload: payload, CreatedAt: time.Now().UTC()}
}

func NewIssueUpdatedMsg(payload any) Message {
	return Message{Type: TypeIssueUpdated, Payload: payload, CreatedAt: time.Now().UTC()}
}

func NewCommentAddedMsg(payload any) Message {
	return Message{Type: TypeCommentAdded, Payload: payload, CreatedAt: time.Now().UTC()}
}

func NewMentionMsg(payload any) Message {
	return Message{Type: TypeMention, Payload: payload, CreatedAt: time.Now().UTC()}
}
