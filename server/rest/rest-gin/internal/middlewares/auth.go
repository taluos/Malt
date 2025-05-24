package middleware

import (
	"github.com/taluos/Malt/server/rest/rest-gin/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthenticMiddleware(authStrategy *auth.AuthOperator) gin.HandlerFunc {
	return authStrategy.AuthFunc()
}
