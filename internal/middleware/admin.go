package middleware

import "github.com/gin-gonic/gin"

// AdminOnly allows only admin users to access the route
func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {

		role := c.GetHeader("Role")

		if role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Admin access required",
			})
			return
		}

		c.Next()
	}
}
