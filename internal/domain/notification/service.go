package notification

import (
	"context"
	"errors"

	"notification/internal/apperror"

	base "notification/internal/repository/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationService struct {
	repo *NotificationRepository
}

func NewNotificationService(repo *NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, n *Notification) error {
	if err := s.repo.Insert(ctx, n); err != nil {
		return apperror.InternalWithDetails("cannot create notification", err)
	}

	err := s.repo.UpdateStatus(ctx, n.Id.Hex(), StatusProcessing, "")
	if err != nil {
		return apperror.InternalWithDetails("Can not Update Notification Status: ", err)
	}
	//TODO: Dispatch notification to channels

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

func (s *NotificationService) ListByUser(ctx context.Context, userId string, limit int, offset int) ([]*Notification, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	total, err := s.repo.Count(ctx, base.WithFilter(bson.M{"user_id": userId}))
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("Cannot count notifications: ", err)
	}
	notifications, err := s.repo.FindAll(ctx,
		base.WithFilter(bson.M{"user_id": userId}),
		base.WithSort("created_at", false),
		base.WithLimit(limit),
		base.WithOffset(offset),
	)
	if err != nil {
		return nil, 0, apperror.InternalWithDetails("Can not get notifications by id: ", err)
	}
	return notifications, total, nil
}
