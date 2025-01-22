package service

import (
	"context"
	"testing"
	"time"

	"middleware-platform/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricsRepository struct {
	mock.Mock
}

func (m *MockMetricsRepository) Create(metrics *model.Metrics) error {
	args := m.Called(metrics)
	return args.Error(0)
}

func (m *MockMetricsRepository) FindLatestByMiddlewareID(middlewareID uint) ([]model.Metrics, error) {
	args := m.Called(middlewareID)
	return args.Get(0).([]model.Metrics), args.Error(1)
}

func (m *MockMetricsRepository) FindByTimeRange(middlewareID uint, start, end time.Time) ([]model.Metrics, error) {
	args := m.Called(middlewareID, start, end)
	return args.Get(0).([]model.Metrics), args.Error(1)
}

func (m *MockMetricsRepository) FindLatestByType(metricType string) ([]model.Metrics, error) {
	args := m.Called(metricType)
	return args.Get(0).([]model.Metrics), args.Error(1)
}

func TestMetricsService_GetLatestMetrics(t *testing.T) {
	mockRepo := new(MockMetricsRepository)
	mockMiddlewareRepo := new(MockMiddlewareRepository)
	service := NewMetricsService(mockRepo, mockMiddlewareRepo)

	metrics := []model.Metrics{
		{
			MiddlewareID: 1,
			Type:        "memory_usage",
			Value:       75.5,
			Unit:        "MB",
			Timestamp:   time.Now(),
		},
	}

	mockRepo.On("FindLatestByMiddlewareID", uint(1)).Return(metrics, nil)

	result, err := service.GetLatestMetrics(1)
	assert.NoError(t, err)
	assert.Contains(t, result, "memory_usage")
	assert.Equal(t, "75.50MB", result["memory_usage"])
	mockRepo.AssertExpectations(t)
}

func TestMetricsService_CollectMetrics(t *testing.T) {
	mockRepo := new(MockMetricsRepository)
	mockMiddlewareRepo := new(MockMiddlewareRepository)
	service := NewMetricsService(mockRepo, mockMiddlewareRepo)

	middlewares := []model.Middleware{
		{
			ID:      1,
			Name:    "test-redis",
			Type:    "redis",
			Version: "6.2",
			Host:    "localhost",
			Port:    "6379",
		},
	}

	mockMiddlewareRepo.On("FindAll").Return(middlewares, nil)
	mockRepo.On("Create", mock.AnythingOfType("*model.Metrics")).Return(nil)

	err := service.CollectMetrics(context.Background())
	assert.NoError(t, err)
	mockMiddlewareRepo.AssertExpectations(t)
}