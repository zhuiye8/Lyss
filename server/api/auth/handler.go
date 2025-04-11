package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/pkg/middleware"
)

// Handler 处理认证相关的HTTP请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建新的认证处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册认证相关的路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		
		// 受保护的路由
		protected := auth.Group("/")
		protected.Use(h.authMiddleware.Authenticate())
		{
			protected.GET("/me", h.GetCurrentUser)
		}
	}
}

// Register 处理用户注册请求
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login 处理用户登录请求
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken 处理令牌刷新请求
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.RefreshToken(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的刷新令牌"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCurrentUser 获取当前认证用户的信息
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取用户信息
	user, err := h.service.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
} 