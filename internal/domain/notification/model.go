package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Channel string
type Status string

const (
	ChannelEmail Channel = "email"
	ChannelSms   Channel = "sms"
	ChannelPush  Channel = "push"
)

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
)

type Notification struct {
	Id        bson.ObjectID     `bson:"_id,omitempty"`
	UserId    string            `bson:"user_id"`
	Channel   Channel           `bson:"channel"`
	Status    Status            `bson:"status"`
	Title     string            `bson:"title"`
	Body      string            `bson:"body"`
	Metadata  map[string]string `bson:"metadata"`
	Error     string            `bson:"error"`
	SentAt    *time.Time        `bson:"sent_at"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

func (n *Notification) GetID() bson.ObjectID {
	return n.Id
}

func (n *Notification) SetID(id bson.ObjectID) {
	n.Id = id
}

func (n *Notification) SetTimestamps(createdAt, updatedAt time.Time) {
	n.CreatedAt = createdAt
	n.UpdatedAt = updatedAt
}
