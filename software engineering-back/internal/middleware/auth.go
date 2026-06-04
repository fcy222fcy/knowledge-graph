package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a stub that checks for token presence.
// TODO: implement real JWT token validation.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"data":    nil,
			})
			c.Abort()
			return
		}
		// Stub: accept any non-empty token, set user_id to 1
		c.Set("user_id", uint(1))
		c.Next()
	}
}
