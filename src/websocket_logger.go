package main

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var (
	// Store all active WebSocket connections
	clients = make(map[*websocket.Conn]bool)
	mutex   sync.Mutex
)

// Broadcast message to all clients
func broadcast(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			// If there's an error, remove the client
			delete(clients, client)
			client.Close()
		}
	}
}

// Custom logger that broadcasts to WebSocket clients
type WebSocketLogger struct{}

func (l *WebSocketLogger) Write(p []byte) (n int, err error) {

	// Send to all websocket clients
	broadcast(p)

	return len(p), nil
}
