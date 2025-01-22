package handler

import (
	"net/http"
	"time"

	"middleware-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	service *service.MetricsService
}

func NewMetricsHandler(service *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{service: service}
}

func (h *MetricsHandler) GetMetricsStatus(c *gin.Context) {
	middlewareID := uint(1) // 从请求参数中获取
	metrics, err := h.service.GetLatestMetrics(middlewareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"metrics": metrics,
	})
}

func (h *MetricsHandler) GetPerformanceMetrics(c *gin.Context) {
	middlewareID := uint(1) // 从请求参数中获取
	duration := 24 * time.Hour // 默认查询最近24小时

	metrics, err := h.service.GetPerformanceMetrics(middlewareID, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"performance": metrics,
	})
}