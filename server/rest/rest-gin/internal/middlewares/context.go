package middleware

import "github.com/gin-gonic/gin"

func Context(key string, value string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
		c.Next()
	}
}

func ContextCopy(targetKey string, sourceKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, exists := c.Get(sourceKey); exists {
			c.Set(targetKey, val)
		}
		c.Next()
	}
}

func ContextFromParam(key string, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Param(paramName)
		if value != "" {
			c.Set(key, value)
		}
		c.Next()
	}
}
