package middleware

import (
	"Malt/pkg/log"
	rbac "Malt/server/rest/internal/RBAC"

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
		err := authenticator.Authenticate(c)
		if err != nil {
			log.Errorf("authenticate error: %v", err)
			c.Abort()
			return
		}
		c.Next()
	}
}
