package dashboard

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhuiye8/Lyss/server/pkg/middleware"
)

// StatisticsResponse 统计数据响应
type StatisticsResponse struct {
	AgentCount        int `json:"agentCount"`
	ConversationCount int `json:"conversationCount"`
	UserCount         int `json:"userCount"`
	TokenUsage        int `json:"tokenUsage"`
}

// UsageData 使用趋势数据
type UsageData struct {
	Date          string `json:"date"`
	Conversations int    `json:"conversations"`
	Tokens        int    `json:"tokens"`
}

// TopAgent 热门智能体数据
type TopAgent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Usage       int     `json:"usage"`
	SuccessRate float64 `json:"successRate"`
}

// RecentActivity 最近活动数据
type RecentActivity struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    string `json:"time"`
	UserID  string `json:"userId"`
}

// Handler 仪表盘处理器
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建仪表盘处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册API路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	dashboardGroup := router.Group("/dashboard")
	{
		// 仪表盘API路由
		dashboardGroup.GET("/statistics", h.authMiddleware.Authenticate(), h.GetStatistics)
		dashboardGroup.GET("/trends", h.authMiddleware.Authenticate(), h.GetUsageTrend)
		dashboardGroup.GET("/top-agents", h.authMiddleware.Authenticate(), h.GetTopAgents)
		dashboardGroup.GET("/activities", h.authMiddleware.Authenticate(), h.GetRecentActivities)
	}
}

// GetStatistics 获取统计数据
func (h *Handler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUsageTrend 获取使用趋势
func (h *Handler) GetUsageTrend(c *gin.Context) {
	// 解析天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}

	data, err := h.service.GetUsageTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取使用趋势失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetTopAgents 获取热门智能体
func (h *Handler) GetTopAgents(c *gin.Context) {
	// 解析限制参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	agents, err := h.service.GetTopAgents(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门智能体失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": agents})
}

// GetRecentActivities 获取最近活动
func (h *Handler) GetRecentActivities(c *gin.Context) {
	// 解析限制参数
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	activities, err := h.service.GetRecentActivities(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最近活动失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": activities})
} 