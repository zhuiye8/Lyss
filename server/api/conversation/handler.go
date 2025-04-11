package conversation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/models"
	"github.com/yourusername/agent-platform/server/pkg/middleware"
	"go.uber.org/zap"
)

// Handler 处理对话相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
	logger         *zap.Logger
}

// NewHandler 创建新的对话处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
		logger:         zap.L().With(zap.String("handler", "conversation")),
	}
}

// RegisterRoutes 注册对话相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// 智能体下的对话
	agentConversations := router.Group("/agents/:agent_id/conversations")
	agentConversations.Use(h.authMiddleware.Authenticate())
	{
		agentConversations.GET("", h.GetConversationsByAgentID)
		agentConversations.POST("", h.CreateConversation)
	}

	// 对话操作
	conversations := router.Group("/conversations")
	conversations.Use(h.authMiddleware.Authenticate())
	{
		conversations.GET("/:id", h.GetConversationByID)
		conversations.DELETE("/:id", h.DeleteConversation)
		conversations.GET("/:conv_id/messages", h.GetMessagesByConversationID)
		conversations.POST("/:conv_id/messages", h.SendMessage)
		conversations.POST("/:conv_id/regenerate", h.RegenerateResponse)
	}

	// 消息操作
	messages := router.Group("/messages")
	messages.Use(h.authMiddleware.Authenticate())
	{
		messages.POST("/:id/feedback", h.ProvideFeedback)
	}
}

// GetConversationsByAgentID 获取智能体的所有对话
func (h *Handler) GetConversationsByAgentID(c *gin.Context) {
	// 从URL获取智能体ID
	agentIDStr := c.Param("agent_id")
	agentID, err := uuid.Parse(agentIDStr)
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
	conversations, err := h.service.GetConversationsByAgentID(agentID, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		h.logger.Error("Failed to get conversations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取对话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// GetConversationByID 根据ID获取对话
func (h *Handler) GetConversationByID(c *gin.Context) {
	// 从URL获取对话ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	conversation, err := h.service.GetConversationByID(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "对话不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此对话"})
			return
		}
		h.logger.Error("Failed to get conversation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取对话失败"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

// CreateConversation 创建新的对话
func (h *Handler) CreateConversation(c *gin.Context) {
	// 从URL获取智能体ID
	agentIDStr := c.Param("agent_id")
	agentID, err := uuid.Parse(agentIDStr)
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
	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置智能体ID
	req.AgentID = agentID

	// 调用服务
	conversation, err := h.service.CreateConversation(req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		h.logger.Error("Failed to create conversation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建对话失败"})
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

// GetMessagesByConversationID 获取对话的所有消息
func (h *Handler) GetMessagesByConversationID(c *gin.Context) {
	// 从URL获取对话ID
	convIDStr := c.Param("conv_id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	messages, err := h.service.GetMessagesByConversationID(convID, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "对话不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此对话"})
			return
		}
		h.logger.Error("Failed to get messages", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage 发送消息到对话
func (h *Handler) SendMessage(c *gin.Context) {
	// 从URL获取对话ID
	convIDStr := c.Param("conv_id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果请求流式响应，目前不支持，返回错误
	if req.Stream {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "流式响应暂未实现"})
		return
	}

	// 调用服务
	response, err := h.service.SendMessage(convID, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "对话不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此对话"})
			return
		}
		h.logger.Error("Failed to send message", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送消息失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegenerateResponse 重新生成AI回复
func (h *Handler) RegenerateResponse(c *gin.Context) {
	// 从URL获取对话ID
	convIDStr := c.Param("conv_id")
	convID, err := uuid.Parse(convIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.RegenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果请求流式响应，目前不支持，返回错误
	if req.Stream {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "流式响应暂未实现"})
		return
	}

	// 调用服务
	response, err := h.service.RegenerateResponse(convID, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "对话不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此对话"})
			return
		}
		h.logger.Error("Failed to regenerate response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重新生成回复失败"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ProvideFeedback 提供消息反馈
func (h *Handler) ProvideFeedback(c *gin.Context) {
	// 从URL获取消息ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 解析请求体
	var req models.MessageFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	err = h.service.ProvideFeedback(id, req, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "消息不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此消息"})
			return
		}
		h.logger.Error("Failed to provide feedback", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提供反馈失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteConversation 删除对话
func (h *Handler) DeleteConversation(c *gin.Context) {
	// 从URL获取对话ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 调用服务
	err = h.service.DeleteConversation(id, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "对话不存在"})
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此对话"})
			return
		}
		h.logger.Error("Failed to delete conversation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除对话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
} 