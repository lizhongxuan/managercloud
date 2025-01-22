package model

import "time"

type Middleware struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Type        string    `json:"type" gorm:"not null"` // 中间件类型：Redis, MySQL, RabbitMQ 等
    Version     string    `json:"version"`
    Status      string    `json:"status"`
    Host        string    `json:"host" gorm:"not null"`
    Port        string    `json:"port" gorm:"not null"`
    Credentials string    `json:"credentials,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
} 