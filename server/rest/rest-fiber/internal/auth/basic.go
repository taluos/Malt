package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal"
)

// BasicStrategy defines Basic authentication strategy.
type BasicStrategy struct {
	compare func(username string, password string) bool
}

var _ AuthStrategy = &BasicStrategy{}

// NewBasicStrategy create basic strategy with compare function.
func NewBasicStrategy(compare func(username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

// AuthFunc defines basic strategy as the gin authentication middleware.
func (b BasicStrategy) AuthFunc() fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := strings.SplitN(c.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			internal.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Drop()

			return errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong.")
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			internal.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Drop()

			return errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong.")
		}

		c.Set(UsernameKey, pair[0])
		c.Next()

		return nil
	}
}
