package routes

import (
	"net/http"

	"github.com/ayushvyasonwork/chatapp1/internal/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/ws", handlers.WebsocketHandler)
}
