package ws

import (
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/google/uuid"
	"sync"
)

type Hub struct {
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

var hub = &Hub{
	clients: make(map[uuid.UUID]*Client),
}

// Register user connection
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.UserID] = client
}

// Unregister on disconnect
func (h *Hub) Unregister(userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

// Send notification to a specific user
func (h *Hub) SendNotification(userID uuid.UUID, message interface{}) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		client.Send(message)
	}
}

// Exported helpers
func Register(c *Client) {
	hub.Register(c)
	logging.GetLogger().Info("client registered: " + c.UserID.String())
}
func Unregister(u uuid.UUID)                        { hub.Unregister(u) }
func SendNotification(u uuid.UUID, msg interface{}) { hub.SendNotification(u, msg) }
