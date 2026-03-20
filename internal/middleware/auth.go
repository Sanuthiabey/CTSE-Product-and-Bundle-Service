package middleware

import "github.com/gin-gonic/gin"

// AuthRequired checks if Authorization header exists
func AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization token required",
			})
			return
		}

		// In real systems we would validate JWT here

		c.Next()
	}
}
