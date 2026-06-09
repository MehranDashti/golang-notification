package provider

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChannelType string
type ProviderName string

const (
	ChannelSms   ChannelType = "sms"
	ChannelEmail ChannelType = "email"
	ChannelPush  ChannelType = "push"
)

const (
	ProviderKavenegar ProviderName = "kavenegar"
	ProviderTwilio    ProviderName = "twilio"
	ProviderSMTP      ProviderName = "smtp"
	ProviderFCM       ProviderName = "fcm"
)

type ProviderConfig struct {
	Id          bson.ObjectID     `bson:"_id,omitempty"`
	Channel     ChannelType       `bson:"channel"`
	Provider    ProviderName      `bson:"provider"`
	Credentials map[string]string `bson:"credentials"`
	IsActive    bool              `bson:"is_active"`
	CreatedAt   time.Time         `bson:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"`
}

func (p *ProviderConfig) GetID() bson.ObjectID {
	return p.Id
}

func (p *ProviderConfig) SetID(id bson.ObjectID) {
	p.Id = id
}

func (p *ProviderConfig) SetTimestamps(createdAt, updatedAt time.Time) {
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt
}
