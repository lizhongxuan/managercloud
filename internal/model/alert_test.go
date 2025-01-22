package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlertRule(t *testing.T) {
	rule := AlertRule{
		Type:      "cpu_usage",
		Target:    "*",
		Threshold: "80",
		Operator:  ">",
		Status:    "enabled",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, "cpu_usage", rule.Type)
	assert.Equal(t, "*", rule.Target)
	assert.Equal(t, "80", rule.Threshold)
	assert.Equal(t, ">", rule.Operator)
	assert.Equal(t, "enabled", rule.Status)
}

func TestAlertHistory(t *testing.T) {
	history := AlertHistory{
		RuleID:    1,
		Message:   "CPU usage exceeded threshold",
		Status:    "triggered",
		CreatedAt: time.Now(),
	}

	assert.Equal(t, uint(1), history.RuleID)
	assert.Equal(t, "CPU usage exceeded threshold", history.Message)
	assert.Equal(t, "triggered", history.Status)
} 