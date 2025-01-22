package repository

import (
	"middleware-platform/internal/model"
	"time"

	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	db.AutoMigrate(&model.AlertRule{}, &model.AlertHistory{})
	return &AlertRepository{db: db}
}

func (r *AlertRepository) CreateRule(rule *model.AlertRule) error {
	return r.db.Create(rule).Error
}

func (r *AlertRepository) UpdateRule(rule *model.AlertRule) error {
	return r.db.Save(rule).Error
}

func (r *AlertRepository) FindAllRules() ([]model.AlertRule, error) {
	var rules []model.AlertRule
	err := r.db.Find(&rules).Error
	return rules, err
}

func (r *AlertRepository) FindEnabledRules() ([]model.AlertRule, error) {
	var rules []model.AlertRule
	err := r.db.Where("status = ?", "enabled").Find(&rules).Error
	return rules, err
}

func (r *AlertRepository) CreateHistory(history *model.AlertHistory) error {
	return r.db.Create(history).Error
}

func (r *AlertRepository) FindHistoryByTimeRange(start, end time.Time) ([]model.AlertHistory, error) {
	var history []model.AlertHistory
	err := r.db.Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}