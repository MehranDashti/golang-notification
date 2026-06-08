package rest

import (
	"net/http"
	"notification/internal/domain/notification"
	"notification/internal/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *notification.NotificationService
}

func NewNotificationHandler(service *notification.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var request dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	n := &notification.Notification{
		UserId:   request.UserID,
		Channel:  notification.Channel(request.Channel),
		Title:    request.Title,
		Body:     request.Body,
		Metadata: request.Metadata,
	}
	if err := h.service.Create(c.Request.Context(), n); err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusCreated,
		"Notification Sent", toResponse(n))
}

func (h *NotificationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	n, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK,
		"Notification Sent", toResponse(n))
}

func (h *NotificationHandler) ListByUser(c *gin.Context) {
	userId := c.Param("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.service.ListByUser(c.Request.Context(), userId, limit, offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := make([]*dto.NotificationResponse, len(notifications))
	for i, n := range notifications {
		resp[i] = toResponse(n)
	}

	dto.RespondSuccess(c, http.StatusOK,
		"Notification Sent", resp)
}

func toResponse(n *notification.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:        n.Id.Hex(),
		UserID:    n.UserId,
		Channel:   string(n.Channel),
		Status:    string(n.Status),
		Title:     n.Title,
		Body:      n.Body,
		Metadata:  n.Metadata,
		SentAt:    n.SentAt,
		CreatedAt: n.CreatedAt,
	}
}
