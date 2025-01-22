package service

import (
	"context"
	"fmt"
	"middleware-platform/internal/model"
	"middleware-platform/internal/repository"
	"strconv"
	"time"
)

type AlertService struct {
	alertRepo   *repository.AlertRepository
	metricsRepo *repository.MetricsRepository
}

func NewAlertService(alertRepo *repository.AlertRepository, metricsRepo *repository.MetricsRepository) *AlertService {
	return &AlertService{
		alertRepo:   alertRepo,
		metricsRepo: metricsRepo,
	}
}

func (s *AlertService) CreateRule(rule *model.AlertRule) error {
	return s.alertRepo.CreateRule(rule)
}

func (s *AlertService) UpdateRule(rule *model.AlertRule) error {
	return s.alertRepo.UpdateRule(rule)
}

func (s *AlertService) GetRules() ([]model.AlertRule, error) {
	return s.alertRepo.FindAllRules()
}

func (s *AlertService) CheckAlerts(ctx context.Context) error {
	rules, err := s.alertRepo.FindEnabledRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		// 获取最新指标
		metrics, err := s.metricsRepo.FindLatestByType(rule.Type)
		if err != nil {
			continue
		}

		// 检查每个指标是否触发告警
		for _, metric := range metrics {
			if s.shouldTriggerAlert(rule, metric) {
				// 创建告警历史记录
				alert := &model.AlertHistory{
					RuleID:  rule.ID,
					Message: fmt.Sprintf("%s exceeded threshold: %v%s (threshold: %s)",
						rule.Type, metric.Value, metric.Unit, rule.Threshold),
					Status: "triggered",
				}
				s.alertRepo.CreateHistory(alert)
			}
		}
	}

	return nil
}

func (s *AlertService) shouldTriggerAlert(rule model.AlertRule, metric model.Metrics) bool {
	threshold, err := strconv.ParseFloat(rule.Threshold, 64)
	if err != nil {
		return false
	}

	switch rule.Operator {
	case ">":
		return metric.Value > threshold
	case ">=":
		return metric.Value >= threshold
	case "<":
		return metric.Value < threshold
	case "<=":
		return metric.Value <= threshold
	case "=":
		return metric.Value == threshold
	default:
		return false
	}
}

func (s *AlertService) GetAlertHistory(startTime, endTime time.Time) ([]model.AlertHistory, error) {
	return s.alertRepo.FindHistoryByTimeRange(startTime, endTime)
}