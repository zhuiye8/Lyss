package config

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/models"
	"github.com/yourusername/agent-platform/server/pkg/middleware"
)

// Handler 处理配置相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建新的配置处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册配置相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	configs := router.Group("/configs")
	{
		// 公开接口
		configs.GET("/system", h.GetSystemConfigs)

		// 需要认证的接口
		protected := configs.Group("/")
		protected.Use(h.authMiddleware.Authenticate())
		{
			protected.POST("", h.UpsertConfig)
			protected.GET("/:id", h.GetConfig)
			protected.DELETE("/:id", h.DeleteConfig)

			protected.GET("/user", h.GetUserConfigs)
			protected.GET("/project/:project_id", h.GetProjectConfigs)
			protected.GET("/application/:application_id", h.GetApplicationConfigs)

			// 管理员接口
			admin := protected.Group("/admin")
			admin.Use(h.authMiddleware.RequireAdmin())
			{
				admin.GET("/all", h.GetAllConfigs)
			}
		}
	}
}

// UpsertConfig 处理创建或更新配置请求
func (h *Handler) UpsertConfig(c *gin.Context) {
	var req UpsertConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	config, err := h.service.UpsertConfig(req, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "设置配置失败"

		if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限设置此配置"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config.ToResponse()})
}

// GetConfig 处理获取单个配置请求
func (h *Handler) GetConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	config, err := h.service.GetConfigByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "获取配置失败"

		if errors.Is(err, ErrConfigNotFound) {
			status = http.StatusNotFound
			errMsg = "配置不存在"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config.ToResponse()})
}

// DeleteConfig 处理删除配置请求
func (h *Handler) DeleteConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	userID, _ := c.Get("user_id")
	err = h.service.DeleteConfig(id, userID.(uuid.UUID))
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "删除配置失败"

		if errors.Is(err, ErrConfigNotFound) {
			status = http.StatusNotFound
			errMsg = "配置不存在"
		} else if errors.Is(err, ErrNoPermission) {
			status = http.StatusForbidden
			errMsg = "没有权限删除此配置"
		}

		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置已删除"})
}

// GetSystemConfigs 处理获取系统配置请求
func (h *Handler) GetSystemConfigs(c *gin.Context) {
	configs, err := h.service.GetSystemConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取系统配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetUserConfigs 处理获取用户配置请求
func (h *Handler) GetUserConfigs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	configs, err := h.service.GetUserConfigs(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetProjectConfigs 处理获取项目配置请求
func (h *Handler) GetProjectConfigs(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	configs, err := h.service.GetProjectConfigs(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetApplicationConfigs 处理获取应用配置请求
func (h *Handler) GetApplicationConfigs(c *gin.Context) {
	applicationID, err := uuid.Parse(c.Param("application_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	configs, err := h.service.GetApplicationConfigs(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取应用配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetAllConfigs 处理获取所有配置请求（仅管理员可用）
func (h *Handler) GetAllConfigs(c *gin.Context) {
	scope := c.Query("scope")
	
	var configs []interface{}
	var err error
	
	if scope == "" {
		// 获取所有配置
		var allConfigs []struct{
			ID        uuid.UUID `json:"id"`
			Key       string    `json:"key"`
			Value     string    `json:"value"`
			Scope     string    `json:"scope"`
			ScopeID   *uuid.UUID `json:"scope_id"`
			CreatedBy *uuid.UUID `json:"created_by"`
			UpdatedBy *uuid.UUID `json:"updated_by"`
			CreatedAt string     `json:"created_at"`
			UpdatedAt string     `json:"updated_at"`
		}
		
		if err := h.service.db.Model(&models.Config{}).Find(&allConfigs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取所有配置失败"})
			return
		}
		
		configs = make([]interface{}, len(allConfigs))
		for i, config := range allConfigs {
			configs[i] = config
		}
	} else {
		// 获取指定作用域的配置
		configs, err = h.service.GetConfigsByScope(scope, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"configs": configs})
}