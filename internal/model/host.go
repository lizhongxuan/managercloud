package model

import (
	"gorm.io/gorm"
)

// Host 主机模型
type Host struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null"`
	IP          string `json:"ip" gorm:"not null"`
	Port        int    `json:"port" gorm:"not null"`
	Username    string `json:"username" gorm:"not null"`
	Password    string `json:"password"`
	SSHKey      string `json:"ssh_key"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// FileSync 文件同步任务模型
type FileSync struct {
	gorm.Model
	HostID        uint    `json:"host_id" gorm:"not null"`
	Host          Host    `json:"host" gorm:"foreignKey:HostID"`
	SourcePath    string  `json:"sourcePath" gorm:"not null"`
	TargetPath    string  `json:"targetPath" gorm:"not null"`
	Status        string  `json:"status"`      // syncing, paused, completed, failed, cancelled
	Progress      float64 `json:"progress"`    // 同步进度 0-100
	Speed         float64 `json:"speed"`       // 传输速度 bytes/s
	LastSyncAt    string  `json:"last_sync_at"`
	Description   string  `json:"description"`
	MD5           string  `json:"md5"`
	ModifiedTime  int64   `json:"modified_time"`
	FileSize      int64   `json:"file_size"`
	IsIncremental bool    `json:"isIncremental"`
	SyncedSize    int64   `json:"synced_size"`  // 已同步大小，用于断点续传
	IsPaused      bool    `json:"is_paused"`    // 是否暂停
}

// FileSyncHistory 文件同步历史记录
type FileSyncHistory struct {
	gorm.Model
	FileSyncID   uint   `json:"file_sync_id" gorm:"not null"`
	FileSync     FileSync `json:"file_sync" gorm:"foreignKey:FileSyncID"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	MD5          string `json:"md5"`
	FileSize     int64  `json:"file_size"`
	SyncType     string `json:"sync_type"` // full: 全量同步, incremental: 增量同步
} 