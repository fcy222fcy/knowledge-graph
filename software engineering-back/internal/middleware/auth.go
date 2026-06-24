package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"software_engineering/pkg/jwt"
)

// RequireAuth JWT 认证中间件，验证请求头中的 Bearer Token 并解析用户信息
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权", "data": nil})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的令牌格式", "data": nil})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效或过期的令牌", "data": nil})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireTeacherAuth 教师权限验证
func RequireTeacherAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行基础认证
		RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		role, exists := c.Get("role")
		if !exists || role.(string) != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "需要教师权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色权限中间件，检查用户是否具有指定角色
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限访问", "data": nil})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "角色信息无效", "data": nil})
			c.Abort()
			return
		}

		// 检查用户角色是否在允许列表中
		for _, allowedRole := range roles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足", "data": nil})
		c.Abort()
	}
}
