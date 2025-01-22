package service

import (
	"middleware-platform/internal/model"
	"middleware-platform/internal/repository"
)

type MiddlewareService struct {
	repo *repository.MiddlewareRepository
}

func NewMiddlewareService(repo *repository.MiddlewareRepository) *MiddlewareService {
	return &MiddlewareService{repo: repo}
}

// GetAll 获取所有中间件列表
func (s *MiddlewareService) GetAll() ([]model.Middleware, error) {
	return s.repo.FindAll()
}

// GetList 根据类型获取中间件列表
func (s *MiddlewareService) GetList(middlewareType string) ([]model.Middleware, error) {
	return s.repo.FindByType(middlewareType)
}

// Create 创建中间件
func (s *MiddlewareService) Create(middleware *model.Middleware) error {
	return s.repo.Create(middleware)
}

// Update 更新中间件
func (s *MiddlewareService) Update(middleware *model.Middleware) error {
	return s.repo.Update(middleware)
}

// Delete 删除中间件
func (s *MiddlewareService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// CheckHealth 检查中间件健康状态
func (s *MiddlewareService) CheckHealth(middleware *model.Middleware) error {
	switch middleware.Type {
	case "redis":
		return s.checkRedisHealth(middleware)
	case "mysql":
		return s.checkMySQLHealth(middleware)
	case "postgresql":
		return s.checkPgHealth(middleware)
	case "zookeeper":
		return s.checkZKHealth(middleware)
	default:
		return nil
	}
} 