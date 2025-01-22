package service

import (
	"context"
	"testing"

	"middleware-platform/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) CreateRule(rule *model.AlertRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockAlertRepository) UpdateRule(rule *model.AlertRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockAlertRepository) FindAllRules() ([]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).([]model.AlertRule), args.Error(1)
}

func (m *MockAlertRepository) FindEnabledRules() ([]model.AlertRule, error) {
	args := m.Called()
	return args.Get(0).([]model.AlertRule), args.Error(1)
}

func (m *MockAlertRepository) CreateHistory(history *model.AlertHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func TestAlertService_CreateRule(t *testing.T) {
	mockRepo := new(MockAlertRepository)
	mockMetricsRepo := new(MockMetricsRepository)
	service := NewAlertService(mockRepo, mockMetricsRepo)

	rule := &model.AlertRule{
		Type:      "cpu_usage",
		Target:    "*",
		Threshold: "80",
		Operator:  ">",
		Status:    "enabled",
	}

	mockRepo.On("CreateRule", rule).Return(nil)

	err := service.CreateRule(rule)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAlertService_CheckAlerts(t *testing.T) {
	mockRepo := new(MockAlertRepository)
	mockMetricsRepo := new(MockMetricsRepository)
	service := NewAlertService(mockRepo, mockMetricsRepo)

	rules := []model.AlertRule{
		{
			ID:        1,
			Type:      "cpu_usage",
			Target:    "*",
			Threshold: "80",
			Operator:  ">",
			Status:    "enabled",
		},
	}

	metrics := []model.Metrics{
		{
			Type:  "cpu_usage",
			Value: 85.0,
			Unit:  "%",
		},
	}

	mockRepo.On("FindEnabledRules").Return(rules, nil)
	mockMetricsRepo.On("FindLatestByType", "cpu_usage").Return(metrics, nil)
	mockRepo.On("CreateHistory", mock.AnythingOfType("*model.AlertHistory")).Return(nil)

	err := service.CheckAlerts(context.Background())
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMetricsRepo.AssertExpectations(t)
} 