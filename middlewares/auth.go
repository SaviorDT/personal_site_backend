package middlewares

import (
	"fmt"
	authController "personal_site/controllers/auth"

	"github.com/gin-gonic/gin"

	"personal_site/schemas"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token_v2")
		if err != nil || token == "" {
			fmt.Println("[Auth Debug] No cookie found:", err)
			c.JSON(401, gin.H{"error": "Authorization cookie is required"})
			c.Abort()
			return
		}

		fmt.Println("[Auth Debug] Cookie found:", token)

		validToken, err := authController.ValidateToken(token)
		if err != nil || !validToken.Valid {
			fmt.Println("[Auth Debug] Token invalid:", err)
			c.JSON(401, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		claims, ok := validToken.Claims.(*schemas.TokenClaims)
		if !ok {
			fmt.Println("[Auth Debug] Invalid claims")
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		user := (&claims.Payload).ExtractUser()

		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token_v2")

		if err != nil || token == "" {
			// No authentication, continue without setting user
			anonymousUser := schemas.TokenUser{
				ID:       0,
				Nickname: "anonymous",
				Role:     "anonymous",
			}

			c.Set("user", anonymousUser)
			// c.Set("user_id", 0) // Optional: might not be needed if endpoints check 'user'
			c.Next()
			return
		}

		validToken, err := authController.ValidateToken(token)
		if err != nil || !validToken.Valid {
			// Invalid token, continue without setting user
			c.JSON(401, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			c.Abort()
			return
		}

		claims, ok := validToken.Claims.(*schemas.TokenClaims)
		if !ok {
			// Invalid claims, continue without setting user
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		user := (&claims.Payload).ExtractUser()
		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}
