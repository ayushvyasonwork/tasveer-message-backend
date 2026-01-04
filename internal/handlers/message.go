package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ayushvyasonwork/chatapp1/internal/db"
	"github.com/ayushvyasonwork/chatapp1/internal/models"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getConversationID(a, b primitive.ObjectID) primitive.ObjectID {
	if a.Hex() < b.Hex() {
		return primitive.NewObjectIDFromTimestamp(a.Timestamp())
	}
	return primitive.NewObjectIDFromTimestamp(b.Timestamp())
}

func handleMessage(senderID primitive.ObjectID, msg IncomingMessage) {
	log.Println("📨 Message received from WS")

	receiverID, err := primitive.ObjectIDFromHex(msg.To)
	if err != nil {
		log.Println("❌ Invalid receiver ID")
		return
	}

	conversationID := getConversationID(senderID, receiverID)

	message := models.Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversationID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Content:        msg.Content,
		MessageType:    "text",
		Status:         "sent",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// ✅ SAVE TO MONGO
	_, err = db.MessageCollection.InsertOne(context.Background(), message)
	if err != nil {
		log.Println("❌ Mongo insert failed:", err)
		return
	}

	log.Println("✅ Message saved to Mongo")

	// 🔎 CHECK IF RECEIVER ONLINE
	clientMu.Lock()
	receiverConn, exists := clients[msg.To]
	clientMu.Unlock()

	if exists {
		message.Status = "delivered"
		message.UpdatedAt = time.Now()

		bytes, _ := json.Marshal(message)
		err := receiverConn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Println("❌ Failed to deliver WS message")
		} else {
			log.Println("📤 Message delivered via WS")
		}
	}

	// 🟢 SEND ACK TO SENDER
	clientMu.Lock()
	senderConn, ok := clients[senderID.Hex()]
	clientMu.Unlock()

	if ok {
		ack := map[string]interface{}{
			"type":    "ack",
			"message": message,
		}
		bytes, _ := json.Marshal(ack)
		senderConn.WriteMessage(websocket.TextMessage, bytes)
	}
}
