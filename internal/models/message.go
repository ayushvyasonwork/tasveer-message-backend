package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	SenderID       primitive.ObjectID `bson:"senderId" json:"senderId"`
	ReceiverID     primitive.ObjectID `bson:"receiverId" json:"receiverId"`
	Content        string             `bson:"content" json:"content"`
	MessageType    string             `bson:"messageType" json:"messageType"`
	MediaURL       string             `bson:"mediaUrl,omitempty" json:"mediaUrl,omitempty"`
	Status         string             `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
