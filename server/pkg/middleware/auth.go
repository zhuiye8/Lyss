package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/agent-platform/server/pkg/auth"
)

// AuthMiddleware 是验证用户认证的中间件
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Authenticate 从请求中提取和验证JWT令牌
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			return
		}

		// 从请求头中提取令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "认证格式无效",
			})
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效或过期的令牌",
			})
			return
		}

		// 将用户信息存入上下文，供后续处理器使用
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole 检查用户是否具有指定角色
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "用户未认证",
			})
			return
		}

		// 检查用户角色是否符合要求
		if userRole != role && userRole != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "权限不足",
			})
			return
		}

		c.Next()
	}
}

// RequireAdmin 检查用户是否具有管理员权限
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
} 