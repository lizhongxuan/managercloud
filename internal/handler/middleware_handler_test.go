package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"middleware-platform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMiddlewareService struct {
	mock.Mock
}

func (m *MockMiddlewareService) GetList(middlewareType string) ([]model.Middleware, error) {
	args := m.Called(middlewareType)
	return args.Get(0).([]model.Middleware), args.Error(1)
}

func (m *MockMiddlewareService) Create(middleware *model.Middleware) error {
	args := m.Called(middleware)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestMiddlewareHandler_GetMiddlewareList(t *testing.T) {
	mockService := new(MockMiddlewareService)
	handler := NewMiddlewareHandler(mockService)
	router := setupTestRouter()
	router.GET("/api/v1/middleware/list", handler.GetMiddlewareList)

	expectedList := []model.Middleware{
		{
			Name:    "test-redis",
			Type:    "redis",
			Version: "6.2",
			Host:    "localhost",
			Port:    "6379",
		},
	}

	mockService.On("GetList", "redis").Return(expectedList, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/middleware/list?type=redis", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []model.Middleware `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, response.Data)
}

func TestMiddlewareHandler_CreateMiddleware(t *testing.T) {
	mockService := new(MockMiddlewareService)
	handler := NewMiddlewareHandler(mockService)
	router := setupTestRouter()
	router.POST("/api/v1/middleware/create", handler.CreateMiddleware)

	mw := model.Middleware{
		Name:    "test-redis",
		Type:    "redis",
		Version: "6.2",
		Host:    "localhost",
		Port:    "6379",
	}

	mockService.On("Create", &mw).Return(nil)

	body, _ := json.Marshal(mw)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/middleware/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
} 