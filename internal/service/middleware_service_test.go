package service

import (
	"testing"

	"middleware-platform/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMiddlewareRepository struct {
	mock.Mock
}

func (m *MockMiddlewareRepository) Create(middleware *model.Middleware) error {
	args := m.Called(middleware)
	return args.Error(0)
}

func (m *MockMiddlewareRepository) FindByType(middlewareType string) ([]model.Middleware, error) {
	args := m.Called(middlewareType)
	return args.Get(0).([]model.Middleware), args.Error(1)
}

func (m *MockMiddlewareRepository) Update(middleware *model.Middleware) error {
	args := m.Called(middleware)
	return args.Error(0)
}

func (m *MockMiddlewareRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMiddlewareRepository) FindAll() ([]model.Middleware, error) {
	args := m.Called()
	return args.Get(0).([]model.Middleware), args.Error(1)
}

func TestMiddlewareService_GetList(t *testing.T) {
	mockRepo := new(MockMiddlewareRepository)
	service := NewMiddlewareService(mockRepo)

	expectedList := []model.Middleware{
		{
			Name:    "test-redis",
			Type:    "redis",
			Version: "6.2",
			Host:    "localhost",
			Port:    "6379",
		},
	}

	mockRepo.On("FindByType", "redis").Return(expectedList, nil)

	list, err := service.GetList("redis")
	assert.NoError(t, err)
	assert.Equal(t, expectedList, list)
	mockRepo.AssertExpectations(t)
}

func TestMiddlewareService_Create(t *testing.T) {
	mockRepo := new(MockMiddlewareRepository)
	service := NewMiddlewareService(mockRepo)

	mw := &model.Middleware{
		Name:    "test-redis",
		Type:    "redis",
		Version: "6.2",
		Host:    "localhost",
		Port:    "6379",
	}

	mockRepo.On("Create", mw).Return(nil)

	err := service.Create(mw)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}