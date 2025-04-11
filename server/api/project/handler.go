package project

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/models"
	"github.com/yourusername/agent-platform/server/pkg/middleware"
)

// Handler 处理项目相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建新的项目处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册项目相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	projects := router.Group("/projects")
	projects.Use(h.authMiddleware.Authenticate())
	{
		projects.POST("", h.CreateProject)
		projects.GET("", h.GetProjects)
		projects.GET("/public", h.GetPublicProjects)
		projects.GET("/:id", h.GetProject)
		projects.PUT("/:id", h.UpdateProject)
		projects.DELETE("/:id", h.DeleteProject)
	}
}

// CreateProject 处理创建项目请求
func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	project, err := h.service.CreateProject(req, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建项目失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project.ToResponse()})
}

// GetProject 处理获取单个项目请求
func (h *Handler) GetProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	userID, _ := c.Get("user_id")
	project, err := h.service.GetProjectByID(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目失败"})
		return
	}

	// 获取项目下的应用数量
	var appsCount int64
	h.service.db.Model(&models.Application{}).Where("project_id = ?", project.ID).Count(&appsCount)
	
	resp := project.ToResponse()
	resp.AppsCount = appsCount
	
	c.JSON(http.StatusOK, gin.H{"project": resp})
}

// UpdateProject 处理更新项目请求
func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	project, err := h.service.UpdateProject(id, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
			return
		}
		if errors.Is(err, ErrNoPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此项目"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新项目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project.ToResponse()})
}

// DeleteProject 处理删除项目请求
func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	userID, _ := c.Get("user_id")
	err = h.service.DeleteProject(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
			return
		}
		if errors.Is(err, ErrNoPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此项目"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除项目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目已删除"})
}

// GetProjects 处理获取项目列表请求
func (h *Handler) GetProjects(c *gin.Context) {
	includeArchived := c.Query("include_archived") == "true"
	userID, _ := c.Get("user_id")
	
	projects, err := h.service.GetProjects(userID.(uuid.UUID), includeArchived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetPublicProjects 处理获取公开项目列表请求
func (h *Handler) GetPublicProjects(c *gin.Context) {
	projects, err := h.service.GetPublicProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公开项目列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
} 