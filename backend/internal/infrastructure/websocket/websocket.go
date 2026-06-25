package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

func newUpgrader() websocket.Upgrader {
	allowed := os.Getenv("FRONTEND_BASE_URL")
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if allowed == "" {
				return true // dev mode: no restriction
			}
			origin := r.Header.Get("Origin")
			return strings.EqualFold(origin, allowed)
		},
	}
}

// Hub manages all active WebSocket connections with room-based broadcasting.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*client]struct{} // userID → connections
	rooms   map[string]map[*client]struct{} // room key → connections
	log     logger.Logger
	done    chan struct{}
}

type client struct {
	userID string
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	rooms  map[string]struct{}
	roomMu sync.Mutex
}

func NewHub(log logger.Logger) *Hub {
	return &Hub{
		clients: make(map[string]map[*client]struct{}),
		rooms:   make(map[string]map[*client]struct{}),
		log:     log,
		done:    make(chan struct{}),
	}
}

// Stop closes all connections cleanly. Call on server shutdown.
func (h *Hub) Stop() {
	close(h.done)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, conns := range h.clients {
		for c := range conns {
			c.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"))
			c.conn.Close()
		}
	}
}

// ServeWS upgrades an HTTP connection to WebSocket and registers the client.
func (h *Hub) ServeWS(ctx context.Context, w http.ResponseWriter, r *http.Request, userID string) {
	upgrader := newUpgrader()
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
		rooms:  make(map[string]struct{}),
	}
	h.register(c)
	defer h.unregister(c)

	go c.writePump(ctx)
	c.readPump(ctx)
}

// BroadcastToRoom sends a message to all clients subscribed to a room.
// Room format: "project:<id>", "issue:<id>", "page:<id>", "space:<id>", "user:<id>".
func (h *Hub) BroadcastToRoom(msg Message) {
	if msg.Room == "" {
		return
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	conns := h.rooms[msg.Room]
	h.mu.RUnlock()

	for c := range conns {
		select {
		case c.send <- data:
		default:
			// slow client — drop rather than block
		}
	}
}

// SendToUser sends a message directly to all connections of a specific user.
func (h *Hub) SendToUser(userID string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	conns := h.clients[userID]
	h.mu.RUnlock()

	for c := range conns {
		select {
		case c.send <- data:
		default:
		}
	}
}

// Online returns true if the user has at least one active connection.
func (h *Hub) Online(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
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
	c.roomMu.Lock()
	rooms := make([]string, 0, len(c.rooms))
	for r := range c.rooms {
		rooms = append(rooms, r)
	}
	c.roomMu.Unlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, room := range rooms {
		if conns, ok := h.rooms[room]; ok {
			delete(conns, c)
			if len(conns) == 0 {
				delete(h.rooms, room)
			}
		}
	}

	conns, ok := h.clients[c.userID]
	if ok {
		if _, exists := conns[c]; exists {
			delete(conns, c)
			close(c.send)
		}
		if len(conns) == 0 {
			delete(h.clients, c.userID)
		}
	}
}

func (h *Hub) subscribeRoom(c *client, room string) {
	h.mu.Lock()
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*client]struct{})
	}
	h.rooms[room][c] = struct{}{}
	h.mu.Unlock()

	c.roomMu.Lock()
	c.rooms[room] = struct{}{}
	c.roomMu.Unlock()
}

func (h *Hub) unsubscribeRoom(c *client, room string) {
	h.mu.Lock()
	if conns, ok := h.rooms[room]; ok {
		delete(conns, c)
		if len(conns) == 0 {
			delete(h.rooms, room)
		}
	}
	h.mu.Unlock()

	c.roomMu.Lock()
	delete(c.rooms, room)
	c.roomMu.Unlock()
}

func (c *client) writePump(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-ctx.Done():
			return
		case <-c.hub.done:
			return
		}
	}
}

func (c *client) readPump(ctx context.Context) {
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var msg Message
		if jsonErr := json.Unmarshal(data, &msg); jsonErr != nil {
			continue
		}
		switch msg.Type {
		case TypeSubscribe:
			if msg.Room != "" {
				c.hub.subscribeRoom(c, msg.Room)
			}
		case TypeUnsubscribe:
			if msg.Room != "" {
				c.hub.unsubscribeRoom(c, msg.Room)
			}
		}
	}
}
