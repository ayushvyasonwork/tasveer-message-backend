package handlers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients  = make(map[string]*websocket.Conn)
	clientMu sync.Mutex
)

func RegisterClient(userID string, conn *websocket.Conn) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clients[userID] = conn
	fmt.Printf("client registered")
}

func RemoveClient(userID string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clients, userID)
}