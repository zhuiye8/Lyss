package agent

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"github.com/zhuiye8/Lyss/server/pkg/middleware"
	"go.uber.org/zap"
)

// Handler 处理智能体相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

// NewHandler 创建新的智能体处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
		logger:         zap.L().With(zap.String("handler", "agent")),
	}
}

// RegisterRoutes 注册智能体相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 应用下的智能体
	applications := router.Group("/applications")
	applications.Use(h.authMiddleware.Authenticate())
	{
		applications.GET("/:app_id/agents", h.GetAgentsByApplicationID)
		applications.POST("/:app_id/agents", h.CreateAgent)
	}

	// 智能体操作
	agents := router.Group("/agents")
	agents.Use(h.authMiddleware.Authenticate())
	{
		agents.GET("/:id", h.GetAgentByID)
		agents.PUT("/:id", h.UpdateAgent)
		agents.DELETE("/:id", h.DeleteAgent)
		agents.PUT("/:id/system-prompt", h.UpdateSystemPrompt)
		agents.PUT("/:id/tools", h.UpdateTools)
		agents.POST("/:id/test", h.TestAgent)
	}
}

// GetAgentsByApplicationID 获取应用下的所有智能体
func (h *Handler) GetAgentsByApplicationID(c *gin.Context) {
	// 从URL获取应用ID
	appIDStr := c.Param("app_id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	agents, err := h.service.GetAgentsByApplicationID(appID, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrApplicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此应用"})
			return
		}
		h.logger.Error("Failed to get agents", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取智能体失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// GetAgentByID 根据ID获取智能体
func (h *Handler) GetAgentByID(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	agent, err := h.service.GetAgentByID(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to get agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取智能体失败"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// CreateAgent 创建新的智能体
func (h *Handler) CreateAgent(c *gin.Context) {
	// 从URL获取应用ID
	appIDStr := c.Param("app_id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的应用ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	agent, err := h.service.CreateAgent(appID, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrApplicationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
			return
		}
		if errors.Is(err, ErrModelConfigNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型配置不存在"})
			return
		}
		if errors.Is(err, ErrKnowledgeBaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "知识库不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此应用"})
			return
		}
		h.logger.Error("Failed to create agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建智能体失败"})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// UpdateAgent 更新智能体
func (h *Handler) UpdateAgent(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	agent, err := h.service.UpdateAgent(id, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrModelConfigNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型配置不存在"})
			return
		}
		if errors.Is(err, ErrKnowledgeBaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "知识库不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to update agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新智能体失败"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// UpdateSystemPrompt 更新智能体系统提示词
func (h *Handler) UpdateSystemPrompt(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.UpdateSystemPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	err = h.service.UpdateSystemPrompt(id, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to update system prompt", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新系统提示词失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateTools 更新智能体工具配置
func (h *Handler) UpdateTools(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.UpdateToolsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	err = h.service.UpdateTools(id, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to update tools", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新工具配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteAgent 删除智能体
func (h *Handler) DeleteAgent(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	err = h.service.DeleteAgent(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to delete agent", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除智能体失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// TestAgent 测试智能体
func (h *Handler) TestAgent(c *gin.Context) {
	// 从URL获取智能体ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的智能体ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.TestAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取智能体信息
	agent, err := h.service.GetAgentByID(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此智能体"})
			return
		}
		h.logger.Error("Failed to get agent for testing", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取智能体失败"})
		return
	}

	// TODO: 当Agent运行时实现后，这里会调用Agent运行
	// 现在先返回一个模拟响应
	mockResponse := map[string]interface{}{
		"id":        uuid.New().String(),
		"agent_id":  agent.ID,
		"query":     req.Message,
		"response":  "这是一个模拟的智能体测试响应。实际功能将在Agent运行时实现后完成",
		"timestamp": c.Request.Context(),
	}

	c.JSON(http.StatusOK, mockResponse)
} 
