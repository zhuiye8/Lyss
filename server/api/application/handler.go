package application

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/pkg/middleware"
)

// Handler 处理应用相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建新的应用处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册应用相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	applications := router.Group("/applications")
	applications.Use(h.authMiddleware.Authenticate())
	{
		applications.POST("", h.CreateApplication)
		applications.GET("/:id", h.GetApplication)
		applications.PUT("/:id", h.UpdateApplication)
		applications.DELETE("/:id", h.DeleteApplication)
		applications.GET("/project/:project_id", h.GetApplicationsByProject)
	}
}

// CreateApplication 处理创建应用请求
func (h *Handler) CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	application, err := h.service.CreateApplication(req, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "创建应用失败"

		if errors.Is(err, ErrProjectNotFound) {
			status = http.StatusNotFound
			errMsg = "项目不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限在此项目下创建应用"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"application": application.ToResponse()})
}

// GetApplication 处理获取单个应用请求
func (h *Handler) GetApplication(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	userID, _ := c.Get("user_id")
	application, err := h.service.GetApplicationByID(id, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "获取应用失败"

		if errors.Is(err, ErrApplicationNotFound) {
			status = http.StatusNotFound
			errMsg = "应用不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限访问此应用"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": application.ToResponse()})
}

// UpdateApplication 处理更新应用请求
func (h *Handler) UpdateApplication(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	var req UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	application, err := h.service.UpdateApplication(id, req, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "更新应用失败"

		if errors.Is(err, ErrApplicationNotFound) {
			status = http.StatusNotFound
			errMsg = "应用不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限修改此应用"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": application.ToResponse()})
}

// DeleteApplication 处理删除应用请求
func (h *Handler) DeleteApplication(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	userID, _ := c.Get("user_id")
	err = h.service.DeleteApplication(id, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "删除应用失败"

		if errors.Is(err, ErrApplicationNotFound) {
			status = http.StatusNotFound
			errMsg = "应用不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限删除此应用"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "应用已删除"})
}

// GetApplicationsByProject 处理获取项目下应用列表请求
func (h *Handler) GetApplicationsByProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	userID, _ := c.Get("user_id")
	applications, err := h.service.GetApplicationsByProject(projectID, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "获取应用列表失败"

		if errors.Is(err, ErrProjectNotFound) {
			status = http.StatusNotFound
			errMsg = "项目不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限访问此项目"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}
