package routes

import (
	"net/http"

	"github.com/ayushvyasonwork/chatapp1/internal/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/health", handlers.HealthHandler)
	// Note: WebSocket doesn't support traditional middleware wrapping
	// Token validation happens inside WebsocketHandler by reading from cookies
	http.HandleFunc("/ws", handlers.WebSocketHandler)
}
