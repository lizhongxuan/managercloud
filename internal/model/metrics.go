package model

import "time"

type Metrics struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MiddlewareID uint      `json:"middleware_id"`
	Type         string    `json:"type"` // cpu_usage, memory_usage, qps, etc
	Value        float64   `json:"value"`
	Unit         string    `json:"unit"` // %, ms, count/s
	Timestamp    time.Time `json:"timestamp"`
} 