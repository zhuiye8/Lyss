package logging

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/pkg/auth"
)

// Handler 处理日志相关的API请求
type Handler struct {
	service *Service
}

// NewHandler 创建新的日志处理器
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 注册API路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	logGroup := router.Group("/logs")
	{
		// 查询日志列表 - 根据查询条件筛选
		logGroup.GET("", auth.RequireAuth(), h.GetLogs)
		
		// 获取特定日志详情
		logGroup.GET("/:id", auth.RequireAuth(), h.GetLogByID)
		
		// 标记错误为已解决
		logGroup.PATCH("/:id/resolve", auth.RequireAuth(), h.MarkErrorAsResolved)
		
		// 获取日志统计信息
		logGroup.GET("/stats", auth.RequireAuth(), h.GetLogStats)
		
		// 获取实时监控数据
		logGroup.GET("/metrics", auth.RequireAuth(), h.GetMetrics)
	}
}

// GetLogs 根据查询条件获取日志列表
func (h *Handler) GetLogs(c *gin.Context) {
	var params LogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的查询参数"})
		return
	}
	
	// 解析日志类型
	logTypeStr := c.DefaultQuery("type", "all")
	var logType LogType
	switch logTypeStr {
	case "api":
		logType = LogTypeAPI
	case "error":
		logType = LogTypeError
	case "model_call":
		logType = LogTypeModelCall
	case "all":
		logType = LogTypeAll
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志类型"})
		return
	}
	
	// 获取当前用户
	userID := auth.GetUserIDFromContext(c)
	
	// 管理员可以查询所有日志，普通用户只能查询自己的
	if !auth.IsAdmin(c) && params.UserID != "" && params.UserID != userID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查询其他用户的日志"})
		return
	}
	
	if params.UserID == "" && !auth.IsAdmin(c) {
		// 非管理员只能查看自己的日志
		params.UserID = userID.String()
	}
	
	logs, totalCount, err := h.service.GetLogs(params, logType)
	if err != nil {
		if err == ErrInvalidLogType {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败"})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"total": totalCount,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// GetLogByID 获取特定日志的详细信息
func (h *Handler) GetLogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
		return
	}
	
	log, err := h.service.GetLogByID(id)
	if err != nil {
		if err == ErrLogNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败"})
		}
		return
	}
	
	// 检查权限：管理员可以查看所有日志，普通用户只能查看自己的
	userID := auth.GetUserIDFromContext(c)
	if !auth.IsAdmin(c) && log.UserID != nil && *log.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查看此日志"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": log})
}

// MarkErrorAsResolved 标记错误日志为已解决
func (h *Handler) MarkErrorAsResolved(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
		return
	}
	
	// 获取当前用户
	userID := auth.GetUserIDFromContext(c)
	
	// 只有管理员可以标记错误为已解决
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有操作权限"})
		return
	}
	
	err = h.service.MarkErrorAsResolved(id, userID)
	if err != nil {
		if err == ErrLogNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "错误已标记为解决"})
}

// GetLogStats 获取日志统计信息
func (h *Handler) GetLogStats(c *gin.Context) {
	// 解析时间范围参数
	startTimeStr := c.DefaultQuery("start_time", "")
	endTimeStr := c.DefaultQuery("end_time", "")
	
	var startTime, endTime time.Time
	var err error
	
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的开始时间格式"})
			return
		}
	}
	
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的结束时间格式"})
			return
		}
	} else {
		endTime = time.Now()
	}
	
	// 默认查询最近24小时
	if startTime.IsZero() {
		startTime = endTime.Add(-24 * time.Hour)
	}
	
	// 只有管理员可以查看统计信息
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查看统计信息"})
		return
	}
	
	stats, err := h.service.GetLogStats(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": stats,
		"meta": gin.H{
			"start_time": startTime,
			"end_time":   endTime,
		},
	})
}

// GetMetrics 获取实时监控数据
func (h *Handler) GetMetrics(c *gin.Context) {
	// 默认返回最近的系统指标
	timeRange := c.DefaultQuery("range", "1h") // 默认1小时
	
	var duration time.Duration
	switch timeRange {
	case "15m":
		duration = 15 * time.Minute
	case "30m":
		duration = 30 * time.Minute
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "12h":
		duration = 12 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	case "7d":
		duration = 7 * 24 * time.Hour
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的时间范围"})
		return
	}
	
	// 只有管理员可以查看监控数据
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查看监控数据"})
		return
	}
	
	// 获取指定时间范围内的系统指标
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	
	// 查询系统指标数据
	var metrics []struct {
		MetricName  string    `json:"metric_name"`
		MetricValue float64   `json:"metric_value"`
		Unit        string    `json:"unit"`
		Tags        string    `json:"tags"`
		CreatedAt   time.Time `json:"created_at"`
	}
	
	if err := h.service.db.Table("system_metrics").
		Select("metric_name, metric_value, unit, tags, created_at").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Order("created_at ASC").
		Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取监控数据失败"})
		return
	}
	
	// 组织数据，按指标名称分组
	metricsMap := make(map[string][]gin.H)
	for _, m := range metrics {
		point := gin.H{
			"value":     m.MetricValue,
			"unit":      m.Unit,
			"timestamp": m.CreatedAt,
		}
		
		if m.Tags != "" {
			var tags map[string]interface{}
			if err := json.Unmarshal([]byte(m.Tags), &tags); err == nil {
				point["tags"] = tags
			}
		}
		
		metricsMap[m.MetricName] = append(metricsMap[m.MetricName], point)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": metricsMap,
		"meta": gin.H{
			"start_time": startTime,
			"end_time":   endTime,
			"range":      timeRange,
		},
	})
} 