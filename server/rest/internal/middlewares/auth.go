package middleware

import (
	auth "github.com/taluos/Malt/core/auth"

	"github.com/gin-gonic/gin"
)

func AuthenticMiddleware(authStrategy auth.AuthOperator) gin.HandlerFunc {
	return authStrategy.AuthFunc()
}
