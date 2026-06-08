package notification

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const CollectionName = "notifications"

type NotificationRepository struct {
	col *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{col: db.Collection(CollectionName)}
}

func (r *NotificationRepository) Create(ctx context.Context, n *Notification) error {
	n.Id = bson.NewObjectID()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Status = StatusPending

	_, err := r.col.InsertOne(ctx, n)

	return err
}

func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*Notification, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var n Notification
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NotificationRepository) FindByUserId(ctx context.Context, userId string, limit int, offset int) ([]*Notification, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.col.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			slog.Error("failed to close cursor", "error", cerr)
		}
	}()

	var result []*Notification
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status Status, errMsg string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"error":      errMsg,
			"updated_at": time.Now(),
		},
	}

	if status == StatusSent {
		update["$set"].(bson.M)["sent_at"] = time.Now()
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, update)

	return err
}

func (r *NotificationRepository) CountByUserId(ctx context.Context, userId string) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{"user_id": userId})
}
