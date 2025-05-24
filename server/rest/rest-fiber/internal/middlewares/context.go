package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/taluos/Malt/pkg/errors"
)

func Context(key string, value string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
		c.Next()
	}
}

func ContextCopy(targetKey string, sourceKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		val := c.Get(sourceKey)
		if val == "" {
			return errors.New("Cannot find source key")
		}
		c.Set(targetKey, val)
		c.Next()
		return nil
	}
}

func ContextFromParam(key string, paramName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		value := c.Params(paramName)
		if value == "" {
			return errors.New("Cannot find param")
		}
		c.Set(key, value)
		c.Next()
		return nil
	}
}
