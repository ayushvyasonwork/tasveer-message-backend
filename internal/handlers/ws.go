package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var jwtSecret = "1234"

type IncomingMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// extract the token string from request query
	// parse token
	// from token extract the claims
	// upgrad to websocket
	// register connection
	// keep connection alive
	log.Println("🔥 GO WS HANDLER HIT")

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("🔥 GO WS HANDLER HIT 1")
	tokenStr := cookie.Value
	log.Println("🔥 GO WS HANDLER HIT 2")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	log.Println("🔥 GO WS HANDLER HIT 3")
	// if err != nil || !token.Valid {
	// 	http.Error(w, "invalid token", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("🔥 GO WS HANDLER HIT 4")
	claims, ok := token.Claims.(jwt.MapClaims)
	log.Println("🔥 GO WS HANDLER HIT 5")
	if !ok {
		http.Error(w, "invalid claims", http.StatusUnauthorized)
		return
	}

	userIDHex, ok := claims["id"].(string)
	log.Println("🔥 GO WS HANDLER HIT 6")
	if !ok {
		http.Error(w, "invalid token payload", http.StatusUnauthorized)
		return
	}
	senderID, err := primitive.ObjectIDFromHex(userIDHex)
	log.Println("🔥 GO WS HANDLER HIT 7")
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println("🔥 GO WS HANDLER HIT 8")
	if err != nil {
		return
	}

	log.Printf("✅ WS CONNECTED: user=%s remote=%s\n",
		userIDHex,
		r.RemoteAddr,
	)
	// fmt.Printf("the conn value is %+v", conn)
	RegisterClient(userIDHex, conn)
	welcome := map[string]string{
		"type":    "system",
		"message": "connected to go websocket",
	}

	bytes, _ := json.Marshal(welcome)
	conn.WriteMessage(websocket.TextMessage, bytes)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WS closed for user:", userIDHex)
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
