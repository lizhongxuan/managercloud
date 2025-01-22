package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	metrics := Metrics{
		MiddlewareID: 1,
		Type:        "memory_usage",
		Value:       75.5,
		Unit:        "MB",
		Timestamp:   time.Now(),
	}

	assert.Equal(t, uint(1), metrics.MiddlewareID)
	assert.Equal(t, "memory_usage", metrics.Type)
	assert.Equal(t, 75.5, metrics.Value)
	assert.Equal(t, "MB", metrics.Unit)
} 