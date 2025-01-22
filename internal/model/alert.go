package model

import "time"

type AlertRule struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"not null"` // cpu_usage, memory_usage, etc
	Target    string    `json:"target" gorm:"not null"` // middleware id or '*' for all
	Threshold string    `json:"threshold" gorm:"not null"` // 阈值
	Operator  string    `json:"operator" gorm:"not null"` // >, <, >=, <=, =
	Status    string    `json:"status"` // enabled, disabled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AlertHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RuleID    uint      `json:"rule_id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // triggered, resolved
	CreatedAt time.Time `json:"created_at"`
}