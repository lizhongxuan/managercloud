package service

import (
	"context"
	"database/sql"
	"fmt"
	"middleware-platform/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func (s *MetricsService) collectRedisMetrics(ctx context.Context, mw model.Middleware) ([]model.Metrics, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", mw.Host, mw.Port),
		Password: mw.Credentials,
		DB:       0,
	})
	defer client.Close()

	info, err := client.Info(ctx).Result()
	if err != nil {
		return nil, err
	}

	// 解析Redis INFO命令的结果
	infoMap := parseRedisInfo(info)

	metrics := []model.Metrics{}

	// 内存使用
	if memoryStr, ok := infoMap["used_memory"]; ok {
		if memoryBytes, err := strconv.ParseFloat(memoryStr, 64); err == nil {
			metrics = append(metrics, model.Metrics{
				MiddlewareID: mw.ID,
				Type:        "memory_usage",
				Value:       memoryBytes / 1024 / 1024, // 转换为MB
				Unit:        "MB",
				Timestamp:   time.Now(),
			})
		}
	}

	// 连接数
	if clientsStr, ok := infoMap["connected_clients"]; ok {
		if clients, err := strconv.ParseFloat(clientsStr, 64); err == nil {
			metrics = append(metrics, model.Metrics{
				MiddlewareID: mw.ID,
				Type:        "connected_clients",
				Value:       clients,
				Unit:        "",
				Timestamp:   time.Now(),
			})
		}
	}

	// 命令执行数
	if commandsStr, ok := infoMap["total_commands_processed"]; ok {
		if commands, err := strconv.ParseFloat(commandsStr, 64); err == nil {
			metrics = append(metrics, model.Metrics{
				MiddlewareID: mw.ID,
				Type:        "commands_processed",
				Value:       commands,
				Unit:        "",
				Timestamp:   time.Now(),
			})
		}
	}

	return metrics, nil
}

// parseRedisInfo 解析Redis INFO命令返回的字符串
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\n")

	for _, line := range lines {
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return result
}

func (s *MetricsService) collectDBMetrics(ctx context.Context, mw model.Middleware) ([]model.Metrics, error) {
	var dsn string
	if mw.Type == "postgresql" {
		dsn = fmt.Sprintf("host=%s port=%s user=postgres password=%s sslmode=disable",
			mw.Host, mw.Port, mw.Credentials)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			"user", mw.Credentials, mw.Host, mw.Port)
	}

	db, err := sql.Open(mw.Type, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	metrics := []model.Metrics{}

	// 获取当前连接数
	var connections float64
	if mw.Type == "postgresql" {
		err = db.QueryRowContext(ctx, 
			"SELECT count(*) FROM pg_stat_activity").Scan(&connections)
	} else {
		err = db.QueryRowContext(ctx, 
			"SELECT COUNT(1) FROM information_schema.processlist").Scan(&connections)
	}
	if err == nil {
		metrics = append(metrics, model.Metrics{
			MiddlewareID: mw.ID,
			Type:        "connections",
			Value:       connections,
			Unit:        "",
			Timestamp:   time.Now(),
		})
	}

	// 获取慢查询数量
	var slowQueries float64
	if mw.Type == "postgresql" {
		err = db.QueryRowContext(ctx, `
			SELECT count(*) 
			FROM pg_stat_activity 
			WHERE state = 'active' 
			AND now() - query_start > interval '1 second'
		`).Scan(&slowQueries)
	} else {
		err = db.QueryRowContext(ctx, `
			SELECT COUNT(1) 
			FROM information_schema.processlist 
			WHERE TIME > 1
		`).Scan(&slowQueries)
	}
	if err == nil {
		metrics = append(metrics, model.Metrics{
			MiddlewareID: mw.ID,
			Type:        "slow_queries",
			Value:       slowQueries,
			Unit:        "",
			Timestamp:   time.Now(),
		})
	}

	return metrics, nil
} 