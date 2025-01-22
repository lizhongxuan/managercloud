package handler

import (
	"encoding/csv"
	"middleware-platform/internal/model"
	"middleware-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MiddlewareHandler struct {
	service *service.MiddlewareService
}

func NewMiddlewareHandler(service *service.MiddlewareService) *MiddlewareHandler {
	return &MiddlewareHandler{service: service}
}

func (h *MiddlewareHandler) GetMiddlewareList(c *gin.Context) {
	// 支持按类型过滤
	middlewareType := c.Query("type")
	var list []model.Middleware
	var err error
	
	if middlewareType != "" {
		list, err = h.service.GetList(middlewareType)
	} else {
		list, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": list,
		"message": "success",
	})
}

func (h *MiddlewareHandler) CreateMiddleware(c *gin.Context) {
	var middleware model.Middleware
	if err := c.ShouldBindJSON(&middleware); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	if err := h.service.Create(&middleware); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": middleware,
		"message": "success",
	})
}

func (h *MiddlewareHandler) UpdateMiddleware(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "invalid id",
		})
		return
	}

	var middleware model.Middleware
	if err := c.ShouldBindJSON(&middleware); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}
	middleware.ID = uint(id)

	if err := h.service.Update(&middleware); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

func (h *MiddlewareHandler) DeleteMiddleware(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "invalid id",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
	})
}

func (h *MiddlewareHandler) ExportMiddlewareList(c *gin.Context) {
	// 支持按类型过滤
	middlewareType := c.Query("type")
	var list []model.Middleware
	var err error
	
	if middlewareType != "" {
		list, err = h.service.GetList(middlewareType)
	} else {
		list, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=middleware_list.csv")

	// 写入CSV头
	writer := csv.NewWriter(c.Writer)
	headers := []string{"ID", "Name", "Type", "Version", "Host", "Port", "Status", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": "Failed to write CSV headers",
		})
		return
	}

	// 写入数据
	for _, item := range list {
		row := []string{
			strconv.FormatUint(uint64(item.ID), 10),
			item.Name,
			item.Type,
			item.Version,
			item.Host,
			item.Port,
			item.Status,
			item.CreatedAt.Format("2006-01-02 15:04:05"),
			item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"message": "Failed to write CSV data",
			})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": "Failed to flush CSV writer",
		})
		return
	}
} 