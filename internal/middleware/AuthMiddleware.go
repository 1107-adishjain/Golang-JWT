package middleware

import (
	helper "github.com/1107-adishjain/golang-jwt/internal/helpers"
	"github.com/gin-gonic/gin"
	"strings"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authentication logic goes here
		// For example, check for a valid JWT token in the request header
		// If authentication fails, you can abort the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header missing"})
			return
		}
		// prefix bearer remove it because the header is Authorization: Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.VerifyJWT(tokenString) //helper function to verify the jwt token
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
