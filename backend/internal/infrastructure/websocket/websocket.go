package websocket

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Hub manages all active WebSocket connections grouped by userID.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*client]struct{}
	log     logger.Logger
}

type client struct {
	userID string
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
}

func NewHub(log logger.Logger) *Hub {
	return &Hub{
		clients: make(map[string]map[*client]struct{}),
		log:     log,
	}
}

// ServeWS upgrades an HTTP connection to WebSocket and registers the client.
func (h *Hub) ServeWS(ctx context.Context, w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(ctx, "ws: upgrade failed", logger.Error(err))
		return
	}

	c := &client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    h,
	}
	h.register(c)

	go c.writePump(ctx)
	c.readPump(ctx)
}

func (h *Hub) register(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.userID] == nil {
		h.clients[c.userID] = make(map[*client]struct{})
	}
	h.clients[c.userID][c] = struct{}{}
}

func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns, ok := h.clients[c.userID]
	if !ok {
		return
	}
	if _, exists := conns[c]; exists {
		delete(conns, c)
		close(c.send)
	}
	if len(conns) == 0 {
		delete(h.clients, c.userID)
	}
}

// Online returns true if the user has at least one active connection.
func (h *Hub) Online(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}

func (c *client) writePump(ctx context.Context) {
	defer c.conn.Close()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}
