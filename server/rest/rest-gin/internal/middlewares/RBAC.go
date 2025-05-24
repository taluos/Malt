package middleware

import (
	rbac "github.com/taluos/Malt/core/RBAC"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
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
		authHeader := c.GetHeader("Authorization")
		tokenString, err := JWT.ParseTokenFromHTTPContext(authHeader)
		err = authenticator.Authenticate(tokenString, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Errorf("authenticate error: %v", err)
			c.Abort()
			return
		}
		c.Next()
	}
}
