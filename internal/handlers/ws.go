package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ayushvyasonwork/chatapp1/internal/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IncomingMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 GO WS HANDLER HIT")

	// 1️⃣ Upgrade first
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WS upgrade failed:", err)
		return
	}

	// 2️⃣ Read JWT from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("❌ Token cookie missing")
		conn.Close()
		return
	}
	tokenStr := cookie.Value
	// 3️⃣ Parse JWT
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return middlewares.GetJWTSecret(), nil
	})

	// if err != nil || !token.Valid {
	// 	log.Println("❌ Invalid JWT:", err)
	// 	conn.Close()
	// 	return
	// }

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("❌ Invalid JWT claims")
		conn.Close()
		return
	}

	userIDHex, ok := claims["id"].(string)
	if !ok {
		log.Println("❌ Invalid user id in token")
		conn.Close()
		return
	}

	senderID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Println("❌ Invalid ObjectID")
		conn.Close()
		return
	}

	// 4️⃣ Register client
	RegisterClient(userIDHex, conn)

	log.Printf("✅ WS CONNECTED: user=%s\n", userIDHex)

	// 5️⃣ Send welcome
	conn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"type":"system","message":"connected to go websocket"}`),
	)

	// 6️⃣ Message loop
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WS closed:", userIDHex)
			RemoveClient(userIDHex)
			conn.Close()
			break
		}

		var msg IncomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("❌ Invalid WS payload")
			continue
		}

		handleMessage(senderID, msg)
	}
}
