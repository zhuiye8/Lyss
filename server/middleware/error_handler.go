package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"server/pkg/errorcode"
	"server/pkg/response"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stackTrace := string(debug.Stack())
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request", c.Request.URL.Path),
					zap.String("stack", stackTrace),
				)

				// 获取请求ID
				requestID, exists := c.Get("requestId")
				if !exists {
					requestID = "unknown"
				}

				// 返回500错误
				c.JSON(http.StatusInternalServerError, response.StandardResponse{
					Success:   false,
					Code:      errorcode.InternalError,
					Message:   errorcode.ErrorMessages[errorcode.InternalError],
					RequestId: requestID.(string),
					Data:      nil,
					Timestamp: response.GetTimestamp(),
				})

				// 中止后续中间件
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()

		// 处理404错误
		if c.Writer.Status() == http.StatusNotFound {
			c.JSON(http.StatusNotFound, response.StandardResponse{
				Success:   false,
				Code:      errorcode.EndpointNotFound,
				Message:   errorcode.ErrorMessages[errorcode.EndpointNotFound],
				RequestId: response.GenRequestId(c),
				Data:      nil,
				Timestamp: response.GetTimestamp(),
			})
		}
	}
} 