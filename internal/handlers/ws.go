package handlers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// --------------------
// WebSocket upgrader
// --------------------
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// IMPORTANT: allow your frontend origin
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "https://tasveer-one.vercel.app" ||
			origin == "https://tasveer.ayushvyas.me"
	},
}

// --------------------
// Incoming WS message
// --------------------
type IncomingMessage struct {
	Type    string `json:"type"`
	To      string `json:"to,omitempty"`
	Content string `json:"content,omitempty"`
}

// --------------------
// JWT secret helper
// --------------------
func getJWTSecret() []byte {
	return []byte("YOUR_JWT_SECRET") // replace with env-based secret
}

// --------------------
// Main WS handler
// --------------------
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 WS HANDSHAKE REQUEST")

	// 1️⃣ Extract short-lived WS token
	wsToken := r.URL.Query().Get("token")
	if wsToken == "" {
		log.Println("❌ WS token missing")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2️⃣ Validate WS token
	token, err := jwt.Parse(wsToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return getJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		log.Println("❌ Invalid WS token:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("❌ Invalid token claims")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3️⃣ Ensure this is a WS-only token
	if claims["type"] != "ws" {
		log.Println("❌ Not a WS token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		log.Println("❌ Invalid user ID in token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 4️⃣ Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WS upgrade failed:", err)
		return
	}

	// 5️⃣ Register client
	RegisterClient(userID, conn)
	log.Printf("✅ WS CONNECTED: user=%s\n", userID)

	// 6️⃣ Send welcome message
	conn.WriteJSON(map[string]string{
		"type":    "system",
		"message": "connected to chat server",
	})

	// 7️⃣ Read loop
	for {
		var msg IncomingMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ WS DISCONNECTED: user=%s\n", userID)
			RemoveClient(userID)
			conn.Close()
			break
		}

		handleMessage(userID, msg)
	}
}
