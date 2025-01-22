package service

import (
	"context"
	"fmt"
	"middleware-platform/internal/model"
	"middleware-platform/internal/repository"
	"time"
)

type MetricsService struct {
	metricsRepo *repository.MetricsRepository
	middlewareRepo *repository.MiddlewareRepository
}

func NewMetricsService(metricsRepo *repository.MetricsRepository, middlewareRepo *repository.MiddlewareRepository) *MetricsService {
	return &MetricsService{
		metricsRepo:     metricsRepo,
		middlewareRepo:  middlewareRepo,
	}
}

func (s *MetricsService) CollectMetrics(ctx context.Context) error {
	// 获取所有中间件
	middlewares, err := s.middlewareRepo.FindAll()
	if err != nil {
		return err
	}

	// 收集每个中间件的指标
	for _, mw := range middlewares {
		metrics := []model.Metrics{}
		
		// 根据中间件类型收集不同的指标
		switch mw.Type {
		case "redis":
			m, err := s.collectRedisMetrics(ctx, mw)
			if err != nil {
				continue
			}
			metrics = append(metrics, m...)
		case "mysql", "postgresql":
			m, err := s.collectDBMetrics(ctx, mw)
			if err != nil {
				continue
			}
			metrics = append(metrics, m...)
		}

		// 保存指标
		for _, m := range metrics {
			if err := s.metricsRepo.Create(&m); err != nil {
				fmt.Printf("Failed to save metrics for middleware %d: %v\n", mw.ID, err)
			}
		}
	}

	return nil
}

func (s *MetricsService) GetLatestMetrics(middlewareID uint) (map[string]interface{}, error) {
	metrics, err := s.metricsRepo.FindLatestByMiddlewareID(middlewareID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, m := range metrics {
		result[m.Type] = fmt.Sprintf("%.2f%s", m.Value, m.Unit)
	}
	return result, nil
}

func (s *MetricsService) GetPerformanceMetrics(middlewareID uint, duration time.Duration) (map[string][]interface{}, error) {
	end := time.Now()
	start := end.Add(-duration)
	
	metrics, err := s.metricsRepo.FindByTimeRange(middlewareID, start, end)
	if err != nil {
		return nil, err
	}

	// 按类型分组
	result := make(map[string][]interface{})
	for _, m := range metrics {
		if _, ok := result[m.Type]; !ok {
			result[m.Type] = make([]interface{}, 0)
		}
		result[m.Type] = append(result[m.Type], map[string]interface{}{
			"timestamp": m.Timestamp,
			"value": m.Value,
			"unit": m.Unit,
		})
	}

	return result, nil
} 