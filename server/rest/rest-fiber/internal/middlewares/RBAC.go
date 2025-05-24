package middleware

import (
	"strings"

	rbac "github.com/taluos/Malt/core/RBAC"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/pkg/log"

	fiber "github.com/gofiber/fiber/v3"
)

func RBACMiddleware(authenticator *rbac.Authenticator) fiber.Handler {
	if authenticator == nil {
		return func(c fiber.Ctx) error {
			log.Errorf("authenticator is nil")
			c.Drop()
			return errors.New("authenticator is nil")
		}
	}
	return func(c fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			log.Errorf("token is empty")
			c.Drop()
			return errors.New("token is empty")
		}
		err := authenticator.Authenticate(getJWTToken(c), c.Route().Path, c.Method())
		if err != nil {
			log.Errorf("authenticate error: %v", err)
			c.Drop()
			return errors.WithCode(code.UserNoAuthority, "user has no authority")
		}
		c.Next()

		return nil
	}
}

func getJWTToken(c fiber.Ctx) string {
	// 从请求头中获取 Authorization 字段
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// 提取 token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	tokenString := parts[1]

	return tokenString
}
