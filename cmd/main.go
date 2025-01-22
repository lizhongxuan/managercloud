package main

import (
	"context"
	"log"
	"middleware-platform/internal/config"
	"middleware-platform/internal/repository"
	"middleware-platform/internal/router"
	"middleware-platform/internal/service"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化仓储层
	middlewareRepo := repository.NewMiddlewareRepository(db)
	metricsRepo := repository.NewMetricsRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	hostRepo := repository.NewHostRepository(db)

	// 初始化服务层
	middlewareService := service.NewMiddlewareService(middlewareRepo)
	metricsService := service.NewMetricsService(metricsRepo, middlewareRepo)
	alertService := service.NewAlertService(alertRepo, metricsRepo)
	hostService := service.NewHostService(hostRepo)

	// 启动监控和告警后台任务
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				if err := metricsService.CollectMetrics(context.Background()); err != nil {
					log.Printf("Failed to collect metrics: %v", err)
				}
				if err := alertService.CheckAlerts(context.Background()); err != nil {
					log.Printf("Failed to check alerts: %v", err)
				}
			}
		}
	}()

	// 初始化路由
	r := router.InitRouter(
		middlewareService,
		metricsService,
		alertService,
		hostService,
	)

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 