package notification

import (
	"context"
	"errors"
	"log/slog"

	"notification/internal/apperror"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationService struct {
	repo *NotificationRepository
}

func NewNotificationService(repo *NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, n *Notification) error {
	if err := s.repo.Create(ctx, n); err != nil {
		slog.Warn("Error for Create Notification")
		return apperror.InternalWithDetails("Can not Create Notification: ", err)
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

func (s *NotificationService) ListByUser(ctx context.Context, userId string, limit int, offset int) ([]*Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	notifications, err := s.repo.FindByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, apperror.InternalWithDetails("Can not get notifications by id: ", err)
	}
	return notifications, nil
}
