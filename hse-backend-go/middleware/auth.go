package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"hse-backend-go/utils"
)

// RequireAuth - cek header "Authorization: Bearer <token>", set userId/roleId/roleName ke context
func RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token not found"})
		return
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		return
	}

	c.Set("userId", claims.UserID)
	c.Set("roleId", claims.RoleID)
	c.Set("roleName", claims.RoleName)
	c.Next()
}

// RequireRole - use after RequireAuth. Blocks the request unless the
// logged-in user's role is one of the given allowed roles.
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName := c.GetString("roleName")
		for _, r := range allowed {
			if roleName == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You don't have permission to perform this action"})
	}
}
