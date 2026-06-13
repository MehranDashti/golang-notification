package notification

import (
	"context"
	"time"

	base "notification/internal/repository/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const CollectionName = "notifications"

type NotificationRepository struct {
	*base.BaseRepository[*Notification]
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		BaseRepository: base.NewBaseRepository[*Notification](db, CollectionName),
	}
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status string, errMsg string) error {
	fields := bson.M{
		"status": status,
		"error":  errMsg,
	}
	if status == string(StatusSent) {
		now := time.Now()
		fields["sent_at"] = now
	}
	return r.UpdateOne(ctx, id, bson.M{"$set": fields})
}
