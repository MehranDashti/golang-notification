package notification

import (
	"context"
	"errors"

	"notification/internal/apperror"
	"notification/internal/kafka"

	base "notification/internal/repository/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationService struct {
	repo     *NotificationRepository
	producer *kafka.Producer
}

func NewNotificationService(repo *NotificationRepository, producer *kafka.Producer) *NotificationService {
	return &NotificationService{repo: repo, producer: producer}
}

func (s *NotificationService) Create(ctx context.Context, n *Notification) error {
	if err := s.repo.Insert(ctx, n); err != nil {
		return apperror.InternalWithDetails("cannot create notification", err)
	}

	if err := s.producer.Publish(ctx, kafka.NotificationMessage{
		NotificationID: n.Id.Hex(),
		UserID:         n.UserId,
		Channel:        string(n.Channel),
		Title:          n.Title,
		Body:           n.Body,
		Metadata:       n.Metadata,
		Attempt:        0,
	}); err != nil {
		_ = s.repo.UpdateStatus(ctx, n.Id.Hex(), string(StatusFailed), "failed to publish to kafka")
		return apperror.InternalWithDetails("cannot publish notification", err)
	}

	return nil
}

func (s *NotificationService) GetById(ctx context.Context, id string) (*Notification, error) {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.NotFound("notification not found")
		}
		return nil, apperror.InternalWithDetails("cannot find notification", err)
	}
	return n, nil
}

func (s *NotificationService) ListByUser(ctx context.Context, userId string, limit, offset int) ([]*Notification, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	total, err := s.repo.Count(ctx, base.WithFilter(bson.M{"user_id": userId}))
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("cannot count notifications", err)
	}

	notifications, err := s.repo.FindAll(ctx,
		base.WithFilter(bson.M{"user_id": userId}),
		base.WithSort("created_at", false),
		base.WithLimit(limit),
		base.WithOffset(offset),
	)
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("cannot list notifications", err)
	}

	return notifications, total, nil
}
