package handler

import (
	"middleware-platform/internal/model"
	"middleware-platform/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"log"
)

type HostHandler struct {
	service *service.HostService
}

func NewHostHandler(service *service.HostService) *HostHandler {
	return &HostHandler{service: service}
}

func (h *HostHandler) GetHostList(c *gin.Context) {
	hosts, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": hosts,
		"message": "success",
	})
}

func (h *HostHandler) CreateHost(c *gin.Context) {
	var host model.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	if err := h.service.Create(&host); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": host,
		"message": "success",
	})
}

func (h *HostHandler) UpdateHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "invalid id",
		})
		return
	}

	var host model.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}
	host.ID = uint(id)

	if err := h.service.Update(&host); err != nil {
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

func (h *HostHandler) DeleteHost(c *gin.Context) {
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

func (h *HostHandler) SyncFile(c *gin.Context) {
	var fileSync model.FileSync
	if err := c.ShouldBindJSON(&fileSync); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": err.Error(),
		})
		return
	}

	log.Printf("Syncing file: %v", fileSync)
	if err := h.service.SyncFile(&fileSync); err != nil {
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

func (h *HostHandler) GetFileSyncs(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("hostId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "invalid host id",
		})
		return
	}

	fileSyncs, err := h.service.GetFileSyncsByHost(uint(hostID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": fileSyncs,
		"message": "success",
	})
} 