package middleware

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zhuiye8/Lyss/server/models"
)

// MetricsCollector 系统指标收集器
type MetricsCollector struct {
	db              *gorm.DB
	mutex           sync.Mutex
	requestCount    int64
	errorCount      int64
	requestDuration int64
	requestsTotal   int64
	errorsTotal     int64
	lastCollected   time.Time
	collectionRate  time.Duration
}

// NewMetricsCollector 创建新的指标收集器
func NewMetricsCollector(db *gorm.DB, collectionRate time.Duration) *MetricsCollector {
	collector := &MetricsCollector{
		db:             db,
		lastCollected:  time.Now(),
		collectionRate: collectionRate,
	}

	// 启动周期性收集
	go collector.periodicCollection()

	return collector
}

// periodicCollection 周期性收集系统指标
func (mc *MetricsCollector) periodicCollection() {
	ticker := time.NewTicker(mc.collectionRate)
	defer ticker.Stop()

	for range ticker.C {
		mc.collectSystemMetrics()
	}
}

// collectSystemMetrics 收集系统指标并保存到数据库
func (mc *MetricsCollector) collectSystemMetrics() {
	now := time.Now()
	
	// 获取Go运行时统计信息
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// 准备通用标签
	tags := map[string]string{
		"host": "localhost", // 在生产环境中应该是实际主机名
	}
	tagsJSON, _ := json.Marshal(tags)
	
	// 获取当前请求计数和错误计数
	mc.mutex.Lock()
	requestCount := mc.requestCount
	errorCount := mc.errorCount
	requestDuration := mc.requestDuration
	
	// 计算每秒请求数和平均响应时间
	elapsed := now.Sub(mc.lastCollected).Seconds()
	requestsPerSecond := float64(requestCount) / elapsed
	var avgRequestDuration float64
	if requestCount > 0 {
		avgRequestDuration = float64(requestDuration) / float64(requestCount)
	}
	
	// 重置计数器
	mc.requestCount = 0
	mc.errorCount = 0
	mc.requestDuration = 0
	mc.lastCollected = now
	mc.mutex.Unlock()
	
	// 创建指标列表
	metrics := []models.SystemMetric{
		{
			MetricName:  "system.memory.alloc",
			MetricValue: float64(mem.Alloc),
			Unit:        "bytes",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.memory.sys",
			MetricValue: float64(mem.Sys),
			Unit:        "bytes",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.memory.heap_alloc",
			MetricValue: float64(mem.HeapAlloc),
			Unit:        "bytes",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.memory.heap_sys",
			MetricValue: float64(mem.HeapSys),
			Unit:        "bytes",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.goroutines",
			MetricValue: float64(runtime.NumGoroutine()),
			Unit:        "count",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.gc.next",
			MetricValue: float64(mem.NextGC),
			Unit:        "bytes",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "system.gc.num",
			MetricValue: float64(mem.NumGC),
			Unit:        "count",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "api.requests_per_second",
			MetricValue: requestsPerSecond,
			Unit:        "count/s",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "api.errors_per_second",
			MetricValue: float64(errorCount) / elapsed,
			Unit:        "count/s",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
		{
			MetricName:  "api.avg_response_time",
			MetricValue: avgRequestDuration,
			Unit:        "ms",
			Tags:        string(tagsJSON),
			CreatedAt:   now,
		},
	}
	
	// 为每个指标生成UUID
	for i := range metrics {
		metrics[i].ID = uuid.New()
	}
	
	// 异步将指标保存到数据库
	go func(metricList []models.SystemMetric) {
		if len(metricList) > 0 {
			if err := mc.db.Create(&metricList).Error; err != nil {
				// 使用 zap 或其他日志记录错误
			}
		}
	}(metrics)
}

// recordRequest 记录请求指标
func (mc *MetricsCollector) recordRequest(duration time.Duration, isError bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.requestCount++
	mc.requestsTotal++
	mc.requestDuration += duration.Milliseconds()
	
	if isError {
		mc.errorCount++
		mc.errorsTotal++
	}
}

// MetricsMiddleware 创建指标收集中间件
func (mc *MetricsCollector) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 计算处理时间
		duration := time.Since(start)
		
		// 判断是否发生错误
		isError := c.Writer.Status() >= 400
		
		// 记录请求
		mc.recordRequest(duration, isError)
	}
} 
