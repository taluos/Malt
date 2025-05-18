package middleware

import "github.com/gin-gonic/gin"

func Context(key string, value string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, c.GetString(value))
		c.Next()
	}
}
