package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

// GetOrgIDFromContext 从上下文中获取组织ID
func GetOrgIDFromContext(c *gin.Context) uuid.UUID {
	orgID, exists := c.Get("organization_id")
	if !exists {
		return uuid.Nil
	}
	return orgID.(uuid.UUID)
}

// IsAdmin 检查当前用户是否是管理员
func IsAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}
	return role.(string) == "admin"
}

// RequireAuth 返回认证中间件函数
// 这是一个过渡函数，应该改为使用AuthMiddleware.Authenticate()
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这是一个空实现，应该被正确的中间件替代
		c.Next()
	}
} 