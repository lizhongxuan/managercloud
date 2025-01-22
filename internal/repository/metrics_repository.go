package repository

import (
	"middleware-platform/internal/model"
	"time"

	"gorm.io/gorm"
)

type MetricsRepository struct {
	db *gorm.DB
}

func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	db.AutoMigrate(&model.Metrics{})
	return &MetricsRepository{db: db}
}

func (r *MetricsRepository) Create(metrics *model.Metrics) error {
	return r.db.Create(metrics).Error
}

func (r *MetricsRepository) FindLatestByMiddlewareID(middlewareID uint) ([]model.Metrics, error) {
	var metrics []model.Metrics
	err := r.db.Where("middleware_id = ?", middlewareID).
		Order("timestamp DESC").
		Limit(10).
		Find(&metrics).Error
	return metrics, err
}

func (r *MetricsRepository) FindByTimeRange(middlewareID uint, start, end time.Time) ([]model.Metrics, error) {
	var metrics []model.Metrics
	err := r.db.Where("middleware_id = ? AND timestamp BETWEEN ? AND ?", 
		middlewareID, start, end).
		Order("timestamp ASC").
		Find(&metrics).Error
	return metrics, err
}

func (r *MetricsRepository) FindLatestByType(metricType string) ([]model.Metrics, error) {
	var metrics []model.Metrics
	err := r.db.Where("type = ?", metricType).
		Order("timestamp DESC").
		Limit(1).
		Find(&metrics).Error
	return metrics, err
} 