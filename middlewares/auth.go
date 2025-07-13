package middlewares

import (
	"github.com/gin-gonic/gin"
	"personal_site/controllers"

	"personal_site/schemas"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		token = token[len("Bearer "):] // Remove "Bearer " prefix

		validToken, err := controllers.ValidateToken(token)
		if err != nil || !validToken.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		claims, ok := validToken.Claims.(*schemas.TokenClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

		user := (&claims.Payload).ExtractUser()

		c.Set("user", user)

		c.Next()
	}
}