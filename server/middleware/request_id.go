package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 为每个请求生成唯一ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取请求ID，如果没有则生成一个新的
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求ID添加到上下文
		c.Set("requestId", requestID)

		// 添加到响应头
		c.Writer.Header().Set("X-Request-ID", requestID)

		// 继续处理请求
		c.Next()
	}
} 