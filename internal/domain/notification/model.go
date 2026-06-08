package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel string
type Status string

const (
	ChannelEmail Channel = "email"
	ChannelSms   Channel = "sms"
	ChannelPush  Channel = "push"
)

const (
	StatusPending     Status = "pending"
	StatusProccessing Status = "proccessing"
	StatusSent        Status = "sent"
	StatusFailed      Status = "failed"
)

type Notification struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    string             `bson:"user_id"`
	Channel   Channel            `bson:"channel"`
	Status    Status             `bson:"status"`
	Title     string             `bson:"string"`
	Body      string             `bson:"body"`
	Metadata  map[string]string  `bson:"meta_title"`
	Error     string             `bson:"error"`
	SentAt    *time.Time         `bson:"sent_at"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
