package rest

import (
	"notification/internal/domain/provider"
)

type ProviderHandler struct {
	service *provider.ProviderConfigService
}

func NewProviderHandler(service *provider.ProviderConfigService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

// func (h *ProviderHandler) List(c *gin.Context) {
// 	userId := c.Param("user_id")
// 	q := dto.ParsePagination(c)

// 	providersConfig, total, err := h.service.List(c.Request.Context(), userId, q.Limit, q.Offset)
// 	if err != nil {
// 		_ = c.Error(err)
// 		return
// 	}

// 	resp := make([]*dto.ProviderConfigResponse, len(providersConfig))
// 	for i, n := range providersConfig {
// 		resp[i] = toResponse(n)
// 	}

// 	dto.RespondSuccess(c, http.StatusOK, "List of All Provider", dto.NewPaginationResponse(resp, total, q))
// }

// func (h *NotificationHandler) Create(c *gin.Context) {
// 	var request dto.CreateNotificationRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		_ = c.Error(err)
// 		return
// 	}

// 	n := &notification.Notification{
// 		UserId:   request.UserID,
// 		Channel:  notification.Channel(request.Channel),
// 		Title:    request.Title,
// 		Body:     request.Body,
// 		Metadata: request.Metadata,
// 	}
// 	if err := h.service.Create(c.Request.Context(), n); err != nil {
// 		_ = c.Error(err)
// 		return
// 	}

// 	dto.RespondSuccess(c, http.StatusCreated,
// 		"Notification Sent", toResponse(n))
// }

// func toResponse(n *notification.Notification) *dto.NotificationResponse {
// 	return &dto.NotificationResponse{
// 		ID:        n.Id.Hex(),
// 		UserID:    n.UserId,
// 		Channel:   string(n.Channel),
// 		Status:    string(n.Status),
// 		Title:     n.Title,
// 		Body:      n.Body,
// 		Metadata:  n.Metadata,
// 		SentAt:    n.SentAt,
// 		CreatedAt: n.CreatedAt,
// 	}
// }
