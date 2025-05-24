package middleware

import (
	fiber "github.com/gofiber/fiber/v3"
	auth "github.com/taluos/Malt/server/rest/rest-fiber/internal/auth"
)

func AuthenticMiddleware(authStrategy *auth.AuthOperator) fiber.Handler {
	return authStrategy.AuthFunc()
}
