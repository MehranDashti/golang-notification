package rest

import (
	"net/http"
	"notification/internal/domain/provider"
	"notification/internal/dto"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	service *provider.ProviderConfigService
}

func NewProviderHandler(service *provider.ProviderConfigService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

func (h *ProviderHandler) List(c *gin.Context) {
	q := dto.ParsePagination(c)

	configs, total, err := h.service.List(c.Request.Context(), q.Limit, q.Offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := make([]*dto.ProviderConfigResponse, len(configs))
	for i, cfg := range configs {
		resp[i] = toProviderResponse(cfg)
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok",
		dto.NewPaginationResponse(resp, total, q))
}

func (h *ProviderHandler) Create(c *gin.Context) {
	var req dto.CreateProviderConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	cfg := &provider.ProviderConfig{
		Channel:     provider.ChannelType(req.Channel),
		Provider:    provider.ProviderName(req.Provider),
		Credentials: req.Credentials,
		IsActive:    req.IsActive,
	}
	if err := h.service.Create(c.Request.Context(), cfg); err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok", toProviderResponse(cfg))
}

func (h *ProviderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	cfg, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok", toProviderResponse(cfg))
}

func (h *ProviderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProviderConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req.Credentials, &req.IsActive); err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok", nil)
}

func (h *ProviderHandler) Activate(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.SetActive(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok", nil)
}

func (h *ProviderHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.SetInActive(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "Ok", nil)
}

func toProviderResponse(cfg *provider.ProviderConfig) *dto.ProviderConfigResponse {
	return &dto.ProviderConfigResponse{
		ID:          cfg.Id.Hex(),
		Channel:     string(cfg.Channel),
		Provider:    string(cfg.Provider),
		Credentials: cfg.Credentials,
		IsActive:    cfg.IsActive,
		CreatedAt:   cfg.CreatedAt,
		UpdatedAt:   cfg.UpdatedAt,
	}
}
