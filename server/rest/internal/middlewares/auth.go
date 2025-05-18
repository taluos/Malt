package middleware

import (
	auth "Malt/server/rest/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthenticMiddleware(authStrategy auth.AuthOperator) gin.HandlerFunc {
	return authStrategy.AuthFunc()
}
