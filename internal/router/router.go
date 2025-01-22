package router

import (
	"middleware-platform/internal/handler"
	"middleware-platform/internal/middleware"
	"middleware-platform/internal/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	middlewareService *service.MiddlewareService,
	metricsService *service.MetricsService,
	alertService *service.AlertService,
	hostService *service.HostService,
) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	// 处理器
	middlewareHandler := handler.NewMiddlewareHandler(middlewareService)
	metricsHandler := handler.NewMetricsHandler(metricsService)
	alertHandler := handler.NewAlertHandler(alertService)
	hostHandler := handler.NewHostHandler(hostService)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 中间件管理
		mw := api.Group("/middleware")
		{
			mw.GET("/list", middlewareHandler.GetMiddlewareList)
			mw.POST("/create", middlewareHandler.CreateMiddleware)
			mw.PUT("/:id", middlewareHandler.UpdateMiddleware)
			mw.DELETE("/:id", middlewareHandler.DeleteMiddleware)
			mw.GET("/export", middlewareHandler.ExportMiddlewareList)
		}

		// 监控指标
		metrics := api.Group("/metrics")
		{
			metrics.GET("/status", metricsHandler.GetMetricsStatus)
			metrics.GET("/performance", metricsHandler.GetPerformanceMetrics)
		}

		// 告警管理
		alerts := api.Group("/alerts")
		{
			alerts.GET("/list", alertHandler.GetAlertsList)
			alerts.POST("/rules", alertHandler.CreateAlertRule)
			alerts.PUT("/rules/:id", alertHandler.UpdateAlertRule)
		}

		// 主机管理
		hosts := api.Group("/hosts")
		{
			hosts.GET("/list", hostHandler.GetHostList)
			hosts.POST("/create", hostHandler.CreateHost)
			hosts.PUT("/:id", hostHandler.UpdateHost)
			hosts.DELETE("/:id", hostHandler.DeleteHost)
			hosts.POST("/sync", hostHandler.SyncFile)
			hosts.GET("/:hostId/syncs", hostHandler.GetFileSyncs)
		}
	}

	return r
} 