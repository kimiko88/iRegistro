package ws

import (
	"encoding/json"
	"fmt"

	"github.com/k/iRegistro/internal/domain"
)

type Broadcaster struct {
	hub *Hub
}

func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{hub: hub}
}

type NotificationMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (b *Broadcaster) NotifyMarkAdded(mark *domain.Mark) {
	msg := NotificationMessage{
		Type:    "MARK_ADDED",
		Payload: mark,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling notification: %v\n", err)
		return
	}

	// Broadcast to all connected clients (Simple approach)
	// In a real app, we would filter by ClassID or UserID logic in Hub
	b.hub.Broadcast <- bytes
}
