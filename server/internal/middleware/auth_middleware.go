package middleware

import (
	"net/http"
	"strings"

	"github.com/bhumika019579/Vouch/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(JWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "inavlid authorization header format",
			})
			return
		}
		tokenString := parts[1]
		userID, err := utils.VerifyJWT(tokenString, JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}
		c.Set("user_id", userID)
		c.Next()

	}
}
