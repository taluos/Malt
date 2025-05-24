package middleware

import (
	"strings"

	rbac "github.com/taluos/Malt/core/RBAC"
	"github.com/taluos/Malt/pkg/log"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(authenticator *rbac.Authenticator) gin.HandlerFunc {
	if authenticator == nil {
		return func(c *gin.Context) {
			log.Errorf("authenticator is nil")
			c.Abort()
			return
		}
	}
	return func(c *gin.Context) {
		//authHeader := c.GetHeader("Authorization")
		//tokenString, err := JWT.ParseTokenFromHTTPContext(authHeader)
		tokenString := getJWTToken(c)
		err := authenticator.Authenticate(tokenString, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Errorf("authenticate error: %v", err)
			c.Abort()
			return
		}
		c.Next()
	}
}

func getJWTToken(c *gin.Context) string {
	// 从请求头中获取 Authorization 字段
	authHeader := c.GetHeader("Authorization")
	// authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// Authorization header not found.
		return ""
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		// Incorrect Authorization header format.
		return ""
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		// JWT not found.
		return ""
	}
	return token
}
