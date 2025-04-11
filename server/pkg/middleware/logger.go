package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/zhuiye8/Lyss/server/models"
)

// responseWriter 是gin.ResponseWriter的自定义包装
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法，记录响应体
func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 返回一个Gin中间件，记录请求和响应信息
func LoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 生成请求ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 记录请求体（可选，需要考虑性能和安全）
		var requestBody []byte
		if shouldLogRequestBody(path) && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装ResponseWriter以记录响应体
		blw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		
		// 获取用户ID（如果已认证）
		var userID *uuid.UUID
		if id, exists := c.Get("user_id"); exists {
			uid := id.(uuid.UUID)
			userID = &uid
		}
		
		// 构建元数据
		metadata := map[string]interface{}{
			"query_params": c.Request.URL.Query(),
		}
		
		// 添加请求体到元数据（如果已记录）
		if len(requestBody) > 0 {
			var reqJSON interface{}
			if err := json.Unmarshal(requestBody, &reqJSON); err == nil {
				metadata["request_body"] = reqJSON
			}
		}
		
		// 根据设置，添加响应体到元数据
		if shouldLogResponseBody(path, statusCode) && blw.body.Len() > 0 {
			var respJSON interface{}
			if err := json.Unmarshal(blw.body.Bytes(), &respJSON); err == nil {
				metadata["response_body"] = respJSON
			}
		}
		
		// 元数据转为JSON字符串
		metadataJSON, _ := json.Marshal(metadata)
		
		// 构建日志消息
		message := ""
		switch {
		case statusCode >= 500:
			message = "服务器错误"
		case statusCode >= 400:
			message = "客户端错误"
		case statusCode >= 300:
			message = "重定向"
		default:
			message = "成功请求"
		}
		message = method + " " + path + " - " + message
		
		// 确定日志级别
		level := models.LogLevelInfo
		if statusCode >= 500 {
			level = models.LogLevelError
		} else if statusCode >= 400 {
			level = models.LogLevelWarn
		}
		
		// 创建API日志记录
		apiLog := models.APILog{
			Log: models.Log{
				Level:     level,
				Category:  models.LogCategoryAPI,
				Message:   message,
				UserID:    userID,
				Metadata:  string(metadataJSON),
				CreatedAt: time.Now(),
			},
			Method:     method,
			Path:       path,
			StatusCode: statusCode,
			IP:         ip,
			UserAgent:  userAgent,
			Duration:   duration,
			RequestID:  requestID,
		}
		
		// 异步保存日志
		go func(log models.APILog) {
			if err := db.Create(&log).Error; err != nil {
				zap.L().Error("Failed to save API log", zap.Error(err))
			}
		}(apiLog)
		
		// 记录错误日志（如果是服务器错误）
		if statusCode >= 500 {
			// 获取错误信息
			errInterface, exists := c.Get("error")
			var errorMsg string
			if exists {
				if err, ok := errInterface.(error); ok {
					errorMsg = err.Error()
				} else {
					errorMsg = "未知错误"
				}
			} else {
				errorMsg = "服务器内部错误"
			}
			
			errorLog := models.ErrorLog{
				Log: models.Log{
					Level:     models.LogLevelError,
					Category:  models.LogCategoryAPI,
					Message:   errorMsg,
					UserID:    userID,
					Metadata:  string(metadataJSON),
					CreatedAt: time.Now(),
				},
				Source:     "API",
				ErrorCode:  "SERVER_ERROR",
				StackTrace: "",
			}
			
			// 异步保存错误日志
			go func(log models.ErrorLog) {
				if err := db.Create(&log).Error; err != nil {
					zap.L().Error("Failed to save error log", zap.Error(err))
				}
			}(errorLog)
		}
	}
}

// shouldLogRequestBody 判断是否应记录请求体
func shouldLogRequestBody(path string) bool {
	// 不记录包含敏感信息的路径，例如登录、注册等
	sensitivePaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}
	
	for _, p := range sensitivePaths {
		if p == path {
			return false
		}
	}
	
	return true
}

// shouldLogResponseBody 判断是否应记录响应体
func shouldLogResponseBody(path string, statusCode int) bool {
	// 只记录错误响应
	if statusCode >= 400 {
		return true
	}
	
	// 不记录特定路径的响应体
	largePaths := []string{
		"/api/v1/logs",
		"/api/v1/metrics",
	}
	
	for _, p := range largePaths {
		if p == path {
			return false
		}
	}
	
	return true
} 
